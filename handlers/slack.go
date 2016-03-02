package handlers

import (
	"fmt"
	"log"

	"github.com/bluele/slack"
	"github.com/sentinel-tools/eventilator/config"
	"github.com/sentinel-tools/eventilator/parser"
)

func PostNotificationEventToSlackChannel(config config.SlackConfig, event parser.NotificationEvent) (err error) {
	api := slack.New(config.Token)
	channel, err := api.FindChannelByName(config.Channel)
	levelColor := "warning"
	switch event.Role {
	case "sentinel":
		levelColor = "danger"
	}
	switch event.Eventname {
	case "+odown":
		levelColor = "danger"
	case "-odown":
		levelColor = "good"
	}
	// switch through event to determine color
	att := slack.Attachment{Color: levelColor, AuthorName: config.AuthorName}
	att.Title = fmt.Sprintf("Sentinel event")
	eventField := slack.AttachmentField{Title: "Event Name", Value: event.Eventname, Short: true}
	podField := slack.AttachmentField{Title: "Pod Name", Value: event.Podname, Short: true}
	roleField := slack.AttachmentField{Title: "Role", Value: event.Role, Short: true}
	reporterField := slack.AttachmentField{Title: "Reporter", Value: config.AuthorSubname, Short: true}
	att.Fields = []*slack.AttachmentField{&eventField, &podField, &roleField, &reporterField}
	if err != nil {
		return (err)
	}
	msg := "Heads up!"
	atts := []*slack.Attachment{&att}
	msgopt := slack.ChatPostMessageOpt{AsUser: false, Attachments: atts}
	log.Printf("[SLACK] MSG=%+v", msg)
	log.Printf("[SLACK] MSGOPT=%+v", msgopt)
	log.Printf("[SLACK] ATT=%+v", att)
	err = api.ChatPostMessage(channel.Id, msg, &msgopt)
	return err
}
