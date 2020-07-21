package main

import (
	"fmt"

	"github.com/ashwanthkumar/slack-go-webhook"
)

func SlackNotification(webhookUrl, username, message, channel string) error {
	payload := slack.Payload{
		Text:      message,
		Username:  username,
		Channel:   channel,
		IconEmoji: ":monkey_face:",
	}
	err := slack.Send(webhookUrl, "", payload)
	if len(err) > 0 {
		fmt.Printf("error: %s\n", err)
	}
	return nil
}
