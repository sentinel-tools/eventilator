// +build custom

package handlers

import (
	"log"

	"github.com/sentinel-tools/eventilator/parser"
)

func GenericHandleEvent(event parser.NotificationEvent) error {
	log.Printf("[GENERIC] Event: %+v", event)
}

func RegisterHandlers() {
	HandlerMap.SetHandler("+sdown", GenericHandleEvent)
	HandlerMap.SetHandler("-sdown", GenericHandleEvent)
	HandlerMap.SetHandler("+switch-master", GenericHandleEvent)
}
