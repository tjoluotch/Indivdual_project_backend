package twillio

import (
	"encoding/json"
	"errors"
	"fmt"
	"mygosource/ind_proj_backend/envar"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func CourseworkSMSMessage(phone_no, cwk_desc, cwk_due_date, first_name *string) error {

	envar.Variables()


	// Set account keys & information
	accountSid := os.Getenv("ACCOUNT_SID")
	authToken := os.Getenv("AUTH_TOKEN")
	smsUrlSection := os.Getenv("SMS_URL_SUBSECTION")
	urlStr := smsUrlSection + accountSid + "/Messages.json"

	twillioNo := os.Getenv("SMS_NUM")

	cwkMessage := "Hi " + *first_name + ",\n Just sending you a Notification that you have a piece of coursework: " + *cwk_desc + ".\n This is due on " + *cwk_due_date +"."

	// Pack up the data for the message
	msgData := url.Values{}
	msgData.Set("To", *phone_no)
	msgData.Set("From", twillioNo)
	msgData.Set("Body", cwkMessage)
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
