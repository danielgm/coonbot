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
	slackApi       *slack.Slack
	commandPattern *regexp.Regexp
	chttp          *http.ServeMux
)

func main() {
	slackToken = os.Getenv("SLACK_TOKEN")
	log.Printf("Using Slack token: %s", slackToken)

	slackApiToken := os.Getenv("SLACK_API_TOKEN")
	log.Printf("Using Slack API token: %s", slackApiToken)

	slackApi = slack.New(slackApiToken)

	commandPattern = regexp.MustCompile(`^\s*#([\w-]+)\s*$`)

	chttp = http.NewServeMux()
	chttp.Handle("/", http.FileServer(http.Dir("./")))

	http.HandleFunc("/", handler)
	log.Println("Waiting for slash command...")
	err := http.ListenAndServe(":"+os.Getenv("PORT"), nil)
	if err != nil {
		panic(err)
	}
}

func handler(res http.ResponseWriter, req *http.Request) {
	log.Printf("req.URL.Path=%s\n", req.URL.Path)
	if req.URL.Path == "/raccoon.png" {
		chttp.ServeHTTP(res, req)
	} else if req.URL.Path == "/hook" && req.Method == "POST" {
		msg := parseRequest(req)
		if msg != nil && msg["token"][0] == slackToken {
			log.Printf("Slash command found! user=%s, channel=%s, text=\"%s\"", msg["user_name"][0], msg["channel_name"][0], msg["text"][0])

			if msg["user_name"][0] != "slackbot" {
				text := msg["text"][0]
				if commandPattern.MatchString(text) {
					channelName := commandPattern.FindStringSubmatch(text)[1]
					log.Printf("Redirecting conversation from %s (%s) to #%s", msg["channel_name"][0], msg["channel_id"][0], channelName)

					sendRedirectMessage(msg["channel_id"][0], channelName)
				} else {
					fmt.Fprintf(res, "{\"text\": \"Usage: /coonbot #channel-name\"}")
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

func sendRedirectMessage(targetChannelId string, channelName string) {
	params := getPostMessageParameters()

	actualChannelId, timestamp, err := slackApi.PostMessage(targetChannelId, fmt.Sprintf(":door::arrow_right: #%s", channelName), params)
	if err != nil {
		log.Printf("Failed to post message to channel: %s\n", err)
		return
	}
	fmt.Printf("Message successfully sent to channel %s at %s", actualChannelId, timestamp)
}

func getPostMessageParameters() slack.PostMessageParameters {
	params := slack.PostMessageParameters{}
	params.Username = "coonbot"
	params.IconURL = "https://coonbot.herokuapp.com/raccoon.png"
	params.LinkNames = 1
	return params
}
