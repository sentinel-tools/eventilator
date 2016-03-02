package config

import (
	"log"
	"os"
)

func GetDefaultReconfigConfig() Rconfig {
	var rc = Rconfig{RedisAddress: "127.0.0.1", RedisPort: 6379, RedisAuth: ""}
	return rc
}

func GetDefaultSensuConfig() SensuConfig {
	var rc = SensuConfig{Hostname: "127.0.0.1", Port: 8000, Enabled: false, Token: ""}
	return rc
}

func GetDefaultEventilatorConfig() Evconfig {
	var rc = Evconfig{RedisAddress: "127.0.0.1", RedisPort: 6379, RedisAuth: "", Sensu: GetDefaultSensuConfig(), Slack: GetDefaultSlackConfig()}
	return rc
}

func GetDefaultSlackConfig() SlackConfig {
	var c = SlackConfig{Token: "",
		Enabled:    false,
		Channel:    "sentinel-events",
		AuthorName: "eventilator",
		Username:   "eventilator",
		TriggerOn:  []string{"+odown", "-odown", "+sdown", "-sdown", "+promoted-slave"},
	}
	hostname, err := os.Hostname()
	if err != nil {
		log.Printf("Unable to get hostname??")
	} else {
		c.AuthorSubname = hostname
	}
	return c
}
