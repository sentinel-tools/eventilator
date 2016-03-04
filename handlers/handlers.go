package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"

	"github.com/sentinel-tools/eventilator/parser"
	"github.com/therealbill/libredis/client"
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

var redisconn *client.Redis

func SetRedisConnection(ip string, port int, auth string) (err error) {
	if auth > "" {
		redisconn, err = client.DialWithConfig(&client.DialConfig{Address: fmt.Sprintf("%s:%d", ip, port), Password: auth})
	} else {
		redisconn, err = client.Dial(ip, port)
	}
	return err
}

func GetRedisConnection() (rc *client.Redis, err error) {
	if redisconn != nil {
		return redisconn, nil
	}
	return rc, fmt.Errorf("Redis connection not initialized!")
}

func GetMyFQDN() (fqdn string, err error) {
	cmd := exec.Command("/bin/hostname", "-f")
	var out bytes.Buffer
	cmd.Stdout = &out
	err = cmd.Run()
	if err != nil {
		log.Printf("GetFDQN Error: '%v'", err)
		return
	}
	fqdn = out.String()
	fqdn = fqdn[:len(fqdn)-1] // removing EOL
	return
}
