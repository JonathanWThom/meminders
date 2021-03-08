package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/sfreiberg/gotwilio"
)

var (
	twilioAccountSID,
	twilioAuthToken,
	twilioFromNumber,
	twilioToNumber string
)

func getenv(name string) string {
	v := os.Getenv(name)
	if v == "" {
		panic("Missing required environment variable: " + name)
	}

	return v
}

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

func main() {
	msg := os.Args[1]
	twilio := gotwilio.NewTwilioClient(twilioAccountSID, twilioAuthToken)

	twilio.SendSMS(twilioFromNumber, twilioToNumber, msg, "", "")
}
