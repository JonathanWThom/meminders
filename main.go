package main

import (
	"log"
	"time"

	"github.com/joho/godotenv"
	"github.com/sfreiberg/gotwilio"
)

var (
	twilioAccountSID,
	twilioAuthToken,
	twilioFromNumber,
	twilioToNumber string
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	twilioAccountSID = getenv("TWILIO_ACCOUNT_SID")
	twilioAuthToken = getenv("TWILIO_AUTH_TOKEN")
	twilioFromNumber = getenv("TWILIO_FROM_NUMBER")
	twilioToNumber = getenv("TWILIO_TO_NUMBER")
}

var reminders = []Reminder{
	{
		Message:   "this is my message",
		Frequency: "Daily",
		Hour:      20,
		Minute:    43,
		Second:    0,
		Zone:      "America/Los_Angeles",
	},
}

type Sender interface {
	SendSMS(string, string, string, string, string, ...*gotwilio.Option) (*gotwilio.SmsResponse, *gotwilio.Exception, error)
}

func main() {
	twilio := gotwilio.NewTwilioClient(twilioAccountSID, twilioAuthToken)

	for tick := range time.Tick(time.Second * 1) {
		for _, reminder := range reminders {
			go func(reminder Reminder) {
				if reminder.MatchesDayAndTime(tick) {
					reminder.SendMessage(twilio, twilioFromNumber, twilioToNumber)
				}
			}(reminder)
		}
	}
}
