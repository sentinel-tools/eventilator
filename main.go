package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/sentinel-tools/eventilator/handlers"
	"github.com/sentinel-tools/eventilator/parser"
)

func main() {
	log.Printf("Consul address: %s", os.Getenv("CONSUL_ADDRESS"))
	path := strings.Split(os.Args[0], "/")
	rand.Seed(time.Now().UnixNano())
	calledAs := path[len(path)-1]
	handlers.SetRedisConnection("127.0.0.1", 6379, "")
	handlers.RegisterHandlers()

	switch calledAs {
	case "reconfigurator":
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
		//rargs := os.Args[2]
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
				h(event)
			}
		}
	}
}
