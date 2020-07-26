package slack

import (
	"fmt"

	"github.com/ashwanthkumar/slack-go-webhook"
)

type Slack struct {
	WebhookUrl string
	Username   string
	Channel    string
}

func (s Slack) Notification(message string) error {
	payload := slack.Payload{
		Text:      message,
		Username:  s.Username,
		Channel:   s.Channel,
		IconEmoji: ":monkey_face:",
	}
	err := slack.Send(s.WebhookUrl, "", payload)
	if len(err) > 0 {
		fmt.Printf("error: %s\n", err)
	}
	return nil
}
