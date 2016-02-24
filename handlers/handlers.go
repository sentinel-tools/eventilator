package handlers

import (
	"fmt"

	"github.com/sentinel-tools/eventilator/parser"
)

type EventHandler func(parser.NotificationEvent) error
type HandlerMapper map[string]EventHandler

func (hm *HandlerMapper) SetHandler(name string, hfunc EventHandler) error {
	(*hm)[name] = hfunc
	return nil
}

func (hm *HandlerMapper) GetHandler(name string) (hfunc EventHandler, err error) {
	hfunc, exists := (*hm)[name]
	if exists {
		return hfunc, nil
	}
	return hfunc, fmt.Errorf("No handler for %s", name)
}

var HandlerMap HandlerMapper

func init() {
	HandlerMap = make(map[string]EventHandler)
}
