// +build !custom

package handlers

import (
	"fmt"
	"log"
	"time"

	"github.com/sentinel-tools/eventilator/parser"
	"github.com/therealbill/libredis/client"
)

type RedisConnection struct {
	client.Redis
}

var redisconn *client.Redis

func SetRedisConnection(ip string, port int, auth string) (err error) {
	if auth == "" {
		redisconn, err = client.Dial(ip, port)
	} else {
		redisconn, err = client.DialWithConfig(&client.DialConfig{Address: fmt.Sprintf("%s:%d", ip, port), Password: auth})
	}
	return
}

func GetRedisConnection() (rc *client.Redis, err error) {
	return redisconn, nil
}

func UpdateRedisStore(event parser.ReconfigurationEvent) (code int, err error) {
	// this is mostly to demonstrate how you can add stuff to do on various events
	// first ensure we are operating on the right event
	now := time.Now()
	//nowstamp := now.Format("2015:03:07:15:04:05")
	rc, err := GetRedisConnection()
	rkey := fmt.Sprintf("failovers:%s:timestamps", event.Podname)
	tstamp := now.Unix()
	rc.SAdd(rkey, fmt.Sprintf("%d", tstamp))
	rc.SAdd("pods-with-failovers", event.Podname)

	rkey = fmt.Sprintf("failovers:%s:log", event.Podname)
	rc.ZAdd(rkey, float64(tstamp), event.NewMasterIP)

	y := fmt.Sprintf("%d", now.Year())
	ym := fmt.Sprintf("%s:%d", y, now.Month())
	ymd := fmt.Sprintf("%s:%d", ym, now.Day())
	ymdh := fmt.Sprintf("%s:%d", ymd, now.Hour())
	ymdhm := fmt.Sprintf("%s:%d", ymdh, now.Minute())
	all := []string{y, ym, ymd, ymdh, ymdhm}
	keybase := fmt.Sprintf("failovers:success:%s:counters", event.Podname)

	for _, k := range all {
		rc.HIncrBy("failovers:aggregated", k, 1)
		rc.HIncrBy(keybase, k, 1)
		rc.SAdd(fmt.Sprintf("failovers:%s", k), event.Podname)
		rc.ZAdd(fmt.Sprintf("failovers:aggregated-by-time:%s", k), float64(tstamp), event.Podname)
		// Now we are going to store a set daily of failovers for the last 60 days
		switch k {
		case ymd:
			rc.SAdd(fmt.Sprintf("pods-with-failovers:%s", k), event.Podname)
			// set to expire 60 days after last entry
			rc.Expire(fmt.Sprintf("pods-with-failovers:%s", k), 5184000)
		}
	}

	log.Printf("Setting new master to %s:%d", event.NewMasterIP, event.NewMasterPort)
	return 0, nil
}
