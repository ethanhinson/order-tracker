package main

import (
	"encoding/json"
	"net/http"
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
	bytes, err := json.Marshal(message)
	if err != nil {
		return false, err
	}
	body := strings.NewReader(string(bytes))
	req, err := http.NewRequest("POST", os.ExpandEnv("https://api.twilio.com/2010-04-01/Accounts/$TWILIO_ACCOUNT_SID/Messages.json"), body)
	if err != nil {
		return false, err
	}
	req.SetBasicAuth(os.Getenv("TWILIO_ACCOUNT_SID"), os.Getenv("TWILIO_AUTH_TOKEN"))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	return true, nil
}

