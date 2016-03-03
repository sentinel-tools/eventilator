package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

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

func postReconfigurationEvent(url string, event parser.ReconfigurationEvent) (resp *http.Response, err error) {
	client := &http.Client{}
	jsonStr, _ := json.Marshal(event)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		return resp, err
	}
	return client.Do(req)
}

func putReconfigurationEvent(url string, event parser.ReconfigurationEvent) (resp *http.Response, err error) {
	client := &http.Client{}
	jsonStr, _ := json.Marshal(event)
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		return resp, err
	}
	return client.Do(req)
}

func postNotificationEvent(url string, event parser.NotificationEvent) (resp *http.Response, err error) {
	client := &http.Client{}
	jsonStr, _ := json.Marshal(event)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		return resp, err
	}
	return client.Do(req)
}

func putNotificationEvent(url string, event parser.NotificationEvent) (resp *http.Response, err error) {
	client := &http.Client{}
	jsonStr, _ := json.Marshal(event)
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		return resp, err
	}
	return client.Do(req)
}
