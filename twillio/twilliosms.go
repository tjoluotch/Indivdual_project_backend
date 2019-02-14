package twillio

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

func SendTwillioMessage(code, phone_no string) error {

	// Set account keys & information
	accountSid := "AC32cd443ee4fc285c6a8d1b30805ae462"
	authToken := "8342021b04ecfd7990cfe31807ab56f4"
	urlStr := "https://api.twilio.com/2010-04-01/Accounts/" + accountSid + "/Messages.json"

	twillioNo := "+447480534149"

	loginMessage := "Thanks for using Studently, please enter this Code: " + code

	// Pack up the data for the message
	msgData := url.Values{}
	msgData.Set("To", phone_no)
	msgData.Set("From", twillioNo)
	msgData.Set("Body", loginMessage)
	msgDataReader := *strings.NewReader(msgData.Encode())

	// Create HTTP request client
	client := &http.Client{}
	req, _ := http.NewRequest("POST", urlStr, &msgDataReader)
	req.SetBasicAuth(accountSid, authToken)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// Make HTTP POST request and return message SID
	resp, _ := client.Do(req)
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		var data map[string]interface{}
		decoder := json.NewDecoder(resp.Body)
		err := decoder.Decode(&data)
		if err == nil {
			fmt.Println(data["sid"])
			return err
		}
	} else {
		fmt.Println(resp.Status)
		err := errors.New("twillio didn't execute the SMS")
		return err
	}
	return nil
}
