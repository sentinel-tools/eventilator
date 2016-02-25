package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/sentinel-tools/eventilator/handlers"
	"github.com/sentinel-tools/eventilator/parser"
)

type rconfig struct {
	RedisAddress string
	RedisPort    int
	RedisAuth    string
}

var rconf rconfig

func getDefaultReconfigConfig() rconfig {
	var rc = rconfig{RedisAddress: "127.0.0.1", RedisPort: 6379, RedisAuth: ""}
	return rc
}

func init() {
	// set up a default config for reconfigurator
	rconf = getDefaultReconfigConfig()
}
func main() {
	rcfile := "/etc/redis/reconfigurator.conf"

	log.Printf("Consul address: %s", os.Getenv("CONSUL_ADDRESS"))
	path := strings.Split(os.Args[0], "/")
	rand.Seed(time.Now().UnixNano())
	calledAs := path[len(path)-1]
	handlers.RegisterHandlers()

	switch calledAs {
	case "reconfigurator":
		raw, err := ioutil.ReadFile(rcfile)
		rcdata := string(raw)
		if err != nil {
			log.Print("Unable to read configfile for reconfigurator. Using default config.")
		} else {
			if _, err := toml.Decode(rcdata, &rconf); err != nil {
				log.Fatalf("Unable to parse configfile for reconfigurator: %+v", err)
			} else {
				log.Print("parsed config")
			}
		}
		handlers.SetRedisConnection(rconf.RedisAddress, rconf.RedisPort, rconf.RedisAuth)

		rargs := os.Args[1:]
		if len(rargs) != 7 {
			log.Printf("Reconfigurator called with an incorrect number of arguments")
		}
		event, err := parser.ParseReconfiguration(rargs)
		if err != nil {
			log.Printf("PARSE ERROR: %v", err)
		} else {
			fmt.Printf("Event: %+v\n", event)
			switch event.Role {
			case "leader":
				log.Print("Calling UpdateRedisStore")
				code, err := handlers.UpdateRedisStore(event)
				if err != nil {
					log.Printf("updateProxtInfo call error: %v", err)
				} else {
					log.Print("Stored event data")
				}
				os.Exit(code)
			case "observer":
				log.Printf("Running on an observer, no action taken")
			}
		}
	case "eventilator":
		eventtype := os.Args[1]
		args := strings.Split(os.Args[2], " ")
		event, err := parser.ParseNotification(eventtype, args)
		if err != nil {
			log.Printf("%v", err)
		} else {
			h, err := handlers.HandlerMap.GetHandler(event.Eventname)
			if err != nil {
				log.Printf("Error: %v", err)
			} else {
				err = h(event)
				if err != nil {
					log.Printf("Error in handler: %v", err)
					os.Exit(1)
				}
			}
		}
	}
}
