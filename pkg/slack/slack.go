package main

import (
	"fmt"

	"github.com/ashwanthkumar/slack-go-webhook"
)

type Slack struct {
	webhookUrl string
	username   string
	channel    string
	message    string
}

func (s *Slack) Notification() error {
	payload := slack.Payload{
		Text:      s.message,
		Username:  s.username,
		Channel:   s.channel,
		IconEmoji: ":monkey_face:",
	}
	err := slack.Send(s.webhookUrl, "", payload)
	if len(err) > 0 {
		fmt.Printf("error: %s\n", err)
	}
	return nil
}
