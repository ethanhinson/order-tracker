package main

import (
	"net/http"
	"net/url"
	"os"
	"strings"
)

type SMSMessage struct {
	To string
	From string
	Body string
}

type SMSMessenger interface {
	Send(message SMSMessage) (bool, error)
}

type TwilioMessenger struct {}
func (s TwilioMessenger) Send(message SMSMessage) (bool, error) {
	msgData := url.Values{}
	msgData.Set("To", message.To)
	msgData.Set("From", message.From)
	msgData.Set("Body", message.Body)
	msgDataReader := *strings.NewReader(msgData.Encode())

	req, err := http.NewRequest("POST", os.ExpandEnv("https://api.twilio.com/2010-04-01/Accounts/$TWILIO_SID/Messages.json"), &msgDataReader)
	if err != nil {
		return false, err
	}
	req.SetBasicAuth(os.Getenv("TWILIO_SID"), os.Getenv("TWILIO_AUTH_TOKEN"))
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	return true, nil
}

