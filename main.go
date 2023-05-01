package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/slack-go/slack"
	"os"
)

// Define the structure of the AWS Console Sign In event
type AWSConsoleSignInEvent struct {
	Version     string        `json:"version"`
	ID          string        `json:"id"`
	DetailType  string        `json:"detail-type"`
	Source      string        `json:"source"`
	Account     string        `json:"account"`
	Time        string        `json:"time"`
	Region      string        `json:"region"`
	Resources   []interface{} `json:"resources"`
	EventDetail EventDetail   `json:"detail"`
}

// Define the structure of the event detail
type EventDetail struct {
	EventVersion        string              `json:"eventVersion"`
	UserIdentity        UserIdentity        `json:"userIdentity"`
	EventTime           string              `json:"eventTime"`
	EventSource         string              `json:"eventSource"`
	EventName           string              `json:"eventName"`
	AWSRegion           string              `json:"awsRegion"`
	SourceIPAddress     string              `json:"sourceIPAddress"`
	UserAgent           string              `json:"userAgent"`
	RequestParameters   interface{}         `json:"requestParameters"`
	ResponseElements    ResponseElements    `json:"responseElements"`
	AdditionalEventData AdditionalEventData `json:"additionalEventData"`
}

// Define the structure of the user identity
type UserIdentity struct {
	Type        string `json:"type"`
	PrincipalID string `json:"principalId"`
	ARN         string `json:"arn"`
	AccountID   string `json:"accountId"`
}

// Define the structure of the response elements
type ResponseElements struct {
	ConsoleLogin string `json:"ConsoleLogin"`
}

// Define the structure of the additional event data
type AdditionalEventData struct {
	LoginTo       string `json:"LoginTo"`
	MobileVersion string `json:"MobileVersion"`
	MFAUsed       string `json:"MFAUsed"`
}

// Create a function to build a Slack TextBlockObject
func buildTextBlockObject(title string, value string) *slack.TextBlockObject {
	return slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("*%s*\n%s", title, value), false, false)
}

// Define the Lambda function handler
func handler(ctx context.Context, event AWSConsoleSignInEvent) {
	// Set Slack authentication token and channel ID
	SLACK_AUTH_TOKEN := slack.New(os.Getenv("SLACK_AUTH_TOKEN"))
	CHANNEL_ID := os.Getenv("CHANNEL_ID")

	// Create Slack message template
	blocks := []slack.Block{
		slack.NewSectionBlock(
			nil,
			[]*slack.TextBlockObject{
				slack.NewTextBlockObject("mrkdwn", ":mega: *AWS Web Console Login*", false, false),
			},
			nil,
		),
		slack.NewSectionBlock(
			nil,
			[]*slack.TextBlockObject{
				buildTextBlockObject("Account ID", event.EventDetail.UserIdentity.AccountID),
				buildTextBlockObject("Type", event.EventDetail.UserIdentity.Type),
				buildTextBlockObject("Source IP", event.EventDetail.SourceIPAddress),
				buildTextBlockObject("Region", event.Region),
				buildTextBlockObject("MFA Used", event.EventDetail.AdditionalEventData.MFAUsed),
				buildTextBlockObject("ARN", event.EventDetail.UserIdentity.ARN),
				buildTextBlockObject("Date", event.Time),
			},
			nil,
		),
	}

	// Send Slack message
	CHANNEL_ID, timestamp, err := SLACK_AUTH_TOKEN.PostMessage(
		CHANNEL_ID,
		slack.MsgOptionBlocks(blocks...),
		slack.MsgOptionAsUser(true), // Add this if you want that the bot would post message as a user, otherwise it will send response using the default slackbot
	)
	// Check for errors and print message if successful
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	fmt.Printf("Message successfully sent to channel %s at %s", CHANNEL_ID, timestamp)
}

// Define the main function which starts the Lambda handler
func main() {
	lambda.Start(handler)
}
