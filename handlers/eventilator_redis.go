package handlers

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/sentinel-tools/eventilator/parser"
)

func PostNotificationEventToRedis(event parser.NotificationEvent) (code int, err error) {
	// this is mostly to demonstrate how you can add stuff to do on various events
	// first ensure we are operating on the right event
	now := time.Now()
	rc, err := GetRedisConnection()
	f, err := os.OpenFile("/var/log/redis/sentinel.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Printf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)
	if err != nil {
		log.Printf("Redis connect error: '%v'", err)
	}
	if rc == nil {
		log.Printf("Redis connect failure: '%v'", err)
	}

	y := fmt.Sprintf("%d", now.Year())
	ym := fmt.Sprintf("%s:%d", y, now.Month())
	ymd := fmt.Sprintf("%s:%d", ym, now.Day())
	ymdh := fmt.Sprintf("%s:%d", ymd, now.Hour())
	ymdhm := fmt.Sprintf("%s:%d", ymdh, now.Minute())
	all := []string{y, ym, ymd, ymdh, ymdhm}
	var keybase string

	switch event.Eventname {
	case "-failover-abort-no-good-slave":
		keybase = fmt.Sprintf("failovers:failures:%s:counters", event.Podname)
		rkey := fmt.Sprintf("failovers:failures:%s:timestamps", event.Podname)
		tstamp := now.Unix()
		rc.SAdd(rkey, fmt.Sprintf("%d", tstamp))
		rc.SAdd("pods-with-failovers", event.Podname)
		rkey = fmt.Sprintf("failovers:%s:faillog", event.Podname)
		rc.ZAdd(rkey, float64(tstamp), event.NewMasterIP)

		keybase = fmt.Sprintf("failovers:failures:%s:counters", event.Eventname)
		rkey = fmt.Sprintf("failovers:failures:%s:timestamps", event.Eventname)
		rc.SAdd(rkey, fmt.Sprintf("%d", tstamp))
		rkey = fmt.Sprintf("failovers:%s:faillog", event.Eventname)
		rc.ZAdd(rkey, float64(tstamp), event.Podname)

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
	case "+sdown", "-sdown", "+odown", "-odown":
		keybase = fmt.Sprintf("sentinel:warnings:%s:%s:counters", event.Eventname, event.Podname)
		rkey := fmt.Sprintf("sentinel:warnings:%s:%s:timestamps", event.Eventname, event.Podname)
		tstamp := now.Unix()
		rc.SAdd(rkey, fmt.Sprintf("%d", tstamp))
		rc.SAdd("pods-with-warnings", event.Podname)
		rkey = fmt.Sprintf("warnings:%s:faillog", event.Eventname)
		rc.ZAdd(rkey, float64(tstamp), event.Podname)

		for _, k := range all {
			rc.HIncrBy("sentinel:warnings:aggregated", k, 1)
			rc.HIncrBy(keybase, k, 1)
			rc.SAdd(fmt.Sprintf("sentinel:warnings:%s", k), event.Podname)
			rc.ZAdd(fmt.Sprintf("sentinel:warnings-by-time:%s", k), float64(tstamp), event.Podname)
			// by warning event
			rc.HIncrBy(fmt.Sprintf("sentinel:warnings:%s:aggregated", event.Eventname), k, 1)
			// Now we are going to store a set daily of failovers for the last 60 days
			// I'd like to make this a config option
			switch k {
			case ymd:
				rc.SAdd(fmt.Sprintf("pods-with-warnings:%s", k), event.Podname)
				// set to expire 60 days after last entry
				rc.Expire(fmt.Sprintf("pods-with-warnings:%s", k), 5184000)
			}
		}
	}

	return 0, nil
}
