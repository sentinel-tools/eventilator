package handlers

import (
	"fmt"
	"log"
	"os"

	"github.com/bluele/slack"
	"github.com/sentinel-tools/eventilator/config"
	"github.com/sentinel-tools/eventilator/parser"
)

func PostNotificationEventToSlackChannel(config config.SlackConfig, event parser.NotificationEvent) (err error) {
	f, err := os.OpenFile("/var/log/redis/sentinel.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Printf("error opening log file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)
	doTrigger := contains(config.TriggerOn, event.Eventname)
	hostname, err := GetMyFQDN()
	if !doTrigger {
		return nil
	}
	api := slack.New(config.Token)
	channel, err := api.FindChannelByName(config.Channel)
	levelColor := "warning"
	// switch through event to determine attachment color
	switch event.Role {
	case "sentinel":
		levelColor = "danger"
	}
	switch event.Eventname {
	case "+odown", "-failover-abort-no-good-slave":
		levelColor = "danger"
	case "-odown", "-sdown":
		levelColor = "good"
	}
	att := slack.Attachment{Color: levelColor, AuthorName: config.AuthorName}
	att.Title = fmt.Sprintf("Sentinel event")
	eventField := slack.AttachmentField{Title: "Event Name", Value: event.Eventname, Short: true}
	podField := slack.AttachmentField{Title: "Pod Name", Value: event.Podname, Short: true}
	roleField := slack.AttachmentField{Title: "Role", Value: event.Role, Short: true}
	reporterField := slack.AttachmentField{Title: "Reporter", Value: hostname, Short: true}
	att.Fields = []*slack.AttachmentField{&eventField, &podField, &roleField, &reporterField}
	if err != nil {
		return (err)
	}
	var msg string
	switch levelColor {
	case "good":
		msg = "Phew, it has recovered."
	case "danger":
		msg = "UHOH! Something is broken."
	case "warning":
		msg = "Heads up, something isn't looking right."
	}
	atts := []*slack.Attachment{&att}
	msgopt := slack.ChatPostMessageOpt{AsUser: false, Attachments: atts, Username: config.Username}
	err = api.ChatPostMessage(channel.Id, msg, &msgopt)
	return err
}
func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
