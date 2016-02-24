// +build !custom

package handlers

import (
	"log"

	"github.com/sentinel-tools/eventilator/parser"
)

// HandleNewSdown will handle +sdown events
func HandleNewSdown(event parser.NotificationEvent) error {
	switch event.Role {
	case "sentinel":
		log.Printf("SDOWN SENTINEL %s", event.Podname)
	case "master":
		log.Printf("SDOWN MASTER %s", event.Podname)
	case "slave":
		log.Printf("SDOWN SLAVE %s", event.Podname)
	}
	return nil
}

// HandleUnSdown will handle +sdown events
func HandleUnSdown(event parser.NotificationEvent) error {
	log.Printf("Got a -sdown event: %+v", event)
	return nil
}

// HandleNewOdown will handle +odown events
func HandleNewOdown(event parser.NotificationEvent) error {
	log.Printf("Got an +odown event: %+v", event)
	return nil
}

// HandleNewOdown will handle +odown events
func HandleUnOdown(event parser.NotificationEvent) error {
	log.Printf("Got an -odown event: %+v", event)
	return nil
}

// HandleNewOdown will handle +odown events
func HandleAbortedFailoverNoGoodSlave(event parser.NotificationEvent) error {
	log.Printf("Got an aborted failover attempt event: %+v", event)
	return nil
}

// HandleSwitchMaster will handle +switch-master events
func HandleSwitchMaster(event parser.NotificationEvent) error {
	log.Printf("Failover for %s completed. New master is %s:%d", event.Podname, event.NewMasterIP, event.NewMasterPort)
	return nil
}

func RegisterHandlers() {
	HandlerMap.SetHandler("+sdown", HandleNewSdown)
	HandlerMap.SetHandler("-sdown", HandleUnSdown)
	HandlerMap.SetHandler("+odown", HandleNewOdown)
	HandlerMap.SetHandler("-odown", HandleUnOdown)
	HandlerMap.SetHandler("+switch-master", HandleSwitchMaster)
	HandlerMap.SetHandler("-failover-abort-no-good-slave", HandleAbortedFailoverNoGoodSlave)
}
