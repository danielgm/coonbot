package main

import (
	"bytes"
	"fmt"
	"github.com/nlopes/slack"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
)

var (
	slackToken     string
	slackApiToken  string
	commandPattern *regexp.Regexp
)

func main() {
	slackToken = os.Getenv("SLACK_TOKEN")
	log.Printf("Using Slack token: %s", slackToken)

	slackApiToken = os.Getenv("SLACK_API_TOKEN")
	log.Printf("Using Slack API token: %s", slackApiToken)

	commandPattern = regexp.MustCompile(`^\s*#([\w-]+)\s*$`)

	http.HandleFunc("/hook", hook)
	log.Println("Waiting for slash command...")
	err := http.ListenAndServe(":"+os.Getenv("PORT"), nil)
	if err != nil {
		panic(err)
	}
}

func hook(res http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		msg := parseRequest(req)
		if msg != nil && msg["token"][0] == slackToken {
			log.Printf("Slash command found! user=%s, channel=%s, text=\"%s\"", msg["user_name"][0], msg["channel_name"][0], msg["text"][0])

			if msg["user_name"][0] != "slackbot" {
				text := msg["text"][0]
				if commandPattern.MatchString(text) {
					channelName := commandPattern.FindStringSubmatch(text)[1]
					fmt.Fprintf(res, "{\"text\": \"Redirecting conversation to #%s\"}", channelName)

					sendRedirect(msg["channel_id"][0], channelName)
				} else {
					fmt.Fprintf(res, "{\"text\": \"Usage: /redirect #channel-name\"}")
				}
			}
		}
	}
}

func parseRequest(req *http.Request) map[string][]string {
	b := new(bytes.Buffer)
	b.ReadFrom(req.Body)
	s := b.String()
	msg, err := url.ParseQuery(s)
	if err != nil {
		log.Printf("Bad webhook request. data=%s", s)
		return nil
	}
	return msg
}

func sendRedirect(channelId string, channelName string) {
	api := slack.New(slackApiToken)
	params := slack.PostMessageParameters{}
	attachment := slack.Attachment{
		Pretext: "some pretext",
		Text:    "some text",
	}
	params.Attachments = []slack.Attachment{attachment}
	channelId, timestamp, err := api.PostMessage(channelName, "Some text", params)
	if err != nil {
		log.Printf("%s\n", err)
		return
	}
	fmt.Printf("Message successfully sent to channel %s at %s", channelId, timestamp)
}
