package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	//"github.com/BurntSushi/toml"
	"github.com/naoina/toml"
	"github.com/sentinel-tools/eventilator/config"
	"github.com/sentinel-tools/eventilator/handlers"
	"github.com/sentinel-tools/eventilator/parser"
)

var (
	rconf config.Rconfig
	econf config.Evconfig
)

func init() {
	// set up a default config for reconfigurator
	rconf = config.GetDefaultReconfigConfig()
	econf = config.GetDefaultEventilatorConfig()
}
func main() {
	rcfile := "/etc/redis/reconfigurator.conf"
	ecfile := "/etc/redis/eventilator.conf"

	f, err := os.OpenFile("/var/log/redis/sentinel.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Printf("error opening log file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)
	path := strings.Split(os.Args[0], "/")
	rand.Seed(time.Now().UnixNano())
	calledAs := path[len(path)-1]
	handlers.RegisterHandlers()

	log.Printf("called as %s", calledAs)
	switch calledAs {
	case "reconfigurator":
		log.Print("client-reconfig-script has been called")
		raw, err := ioutil.ReadFile(rcfile)
		if err != nil {
			log.Print("Unable to read configfile for reconfigurator. Using default config.")
		} else {
			if err := toml.Unmarshal(raw, &rconf); err != nil {
				log.Fatalf("Unable to parse configfile for reconfigurator: %+v", err)
			} else {
				log.Print("parsed reconfigurator config")
			}
		}
		err = handlers.SetRedisConnection(rconf.RedisAddress, rconf.RedisPort, rconf.RedisAuth)
		if err != nil {
			log.Fatalf("Unable to set up Redis connection. Error='%v'", err)
		}

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
		raw, err := ioutil.ReadFile(ecfile)
		if err != nil {
			log.Print("Unable to read configfile for reconfigurator. Using default config.")
		} else {
			if err := toml.Unmarshal(raw, &econf); err != nil {
				log.Fatalf("Unable to parse configfile for eventilator: %+v", err)
			} else {
				log.Print("parsed eventilator config")
			}
		}
		err = handlers.SetRedisConnection(econf.RedisAddress, econf.RedisPort, econf.RedisAuth)
		if err != nil {
			log.Fatalf("Unable to connect to Store. Error='%v'", err)
		}
		eventtype := os.Args[1]
		args := strings.Split(os.Args[2], " ")
		event, err := parser.ParseNotification(eventtype, args)
		if eventtype == "+new-epoch" {
			os.Exit(0)
		}
		var errors []error
		if err != nil {
			log.Printf("Parse error: %v", err)
			os.Exit(2) // don't retry because we can't parse it anyway
		} else {
			h, err := handlers.HandlerMap.GetHandler(event.Eventname)
			if err != nil {
				if strings.Contains(err.Error(), "No handler") {
				} else {
					log.Printf("GetHandler Error: %v", err)
				}
			} else {
				err = h(event)
				if err != nil {
					log.Printf("Error in handler: %v", err)
					errors = append(errors, err)
				}
			}
			if econf.RedisEnabled {
				_, err = handlers.PostNotificationEventToRedis(event)
				if err != nil {
					log.Printf("Error in Slack handler: %v", err)
				}
			}
			if econf.Slack.Enabled {
				err = handlers.PostNotificationEventToSlackChannel(econf.Slack, event)
				if err != nil {
					log.Printf("Error in Slack handler: %v", err)
				}
			}
			if econf.SensuJIT.Enabled {
				err = handlers.PostNotificationEventToSensuJIT(econf.SensuJIT, event)
				if err != nil {
					log.Printf("Error in SensuJIT handler: %v", err)
				}
			}

		}
		if len(errors) > 0 {
			for _, err := range errors {
				log.Printf("[end]Handler ERROR: %v", err)
			}
			os.Exit(1)
		}
	}
	os.Exit(0)
}
