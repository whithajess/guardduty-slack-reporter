package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/nlopes/slack"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// GuardDutyFinding - Takes relevant info out of the json payload from GuardDuty
type GuardDutyFinding struct {
	Title       string      `json:"title"`
	Description string      `json:"description"`
	Severity    json.Number `json:"severity"`
	Type        string      `json:"type"`
	AccountId   string      `json:"accountId"`
}

// Initialise Slack API with the Bot Token
var api = slack.New(os.Getenv("OAUTH_ACCESS_TOKEN"))

// Reporter - Listens for CloudWatch events of GuardDuty Findings
// Then formats these and sends them to Slack
func Reporter(event events.CloudWatchEvent) (events.CloudWatchEvent, error) {
	var finding GuardDutyFinding
	json.Unmarshal([]byte(event.Detail), &finding)

	params := slack.PostMessageParameters{}
	attachment := slack.Attachment{
		Title: finding.Title,
		Text:  finding.Description,
		Color: "danger",
		Fields: []slack.AttachmentField{
			{
				Title: "Account ID",
				Value: finding.AccountId,
			},
			{
				Title: "Severity",
				Value: string(finding.Severity),
			},
			{
				Title: "Type",
				Value: finding.Type,
			},
		},
	}
	params.Attachments = []slack.Attachment{attachment}
	channelID, timestamp, err := api.PostMessage(os.Getenv("CHANNEL"), "", params)

	// Logging for errors / success
	if err != nil {
		fmt.Printf("%s\n", err)
		fmt.Printf(err.Error())
	} else {
		fmt.Printf("Message successfully sent to channel %s at %s", channelID, timestamp)
	}

	// TODO: Not sure what we are meant to send back to CloudWatch - seems to work
	return event, err
}

func main() {
	lambda.Start(Reporter)
}
