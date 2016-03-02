// +build !custom

package handlers

import (
	"log"

	"github.com/sentinel-tools/eventilator/parser"
)

func sendSensuAlert(event parser.NotificationEvent) {
	log.Print("sendSensuAlert called")
}
