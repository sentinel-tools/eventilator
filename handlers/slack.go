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
	// switch through event to determine color
	att := slack.Attachment{Color: levelColor, AuthorName: config.AuthorName, AuthorSubname: config.AuthorSubname}
	att.Title = fmt.Sprintf("Sentinel event")
	eventField := slack.AttachmentField{Title: "Event Name", Value: event.Eventname, Short: true}
	podField := slack.AttachmentField{Title: "Pod Name", Value: event.Podname, Short: true}
	att.Fields = []*slack.AttachmentField{&eventField, &podField}
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
