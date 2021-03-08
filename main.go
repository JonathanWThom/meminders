package main

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/sfreiberg/gotwilio"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	log "github.com/sirupsen/logrus"
)

var (
	databaseURL,
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

	databaseURL = getenv("DATABASE_URL")
	twilioAccountSID = getenv("TWILIO_ACCOUNT_SID")
	twilioAuthToken = getenv("TWILIO_AUTH_TOKEN")
	twilioFromNumber = getenv("TWILIO_FROM_NUMBER")
	twilioToNumber = getenv("TWILIO_TO_NUMBER")

	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
}

type Sender interface {
	SendSMS(string, string, string, string, string, ...*gotwilio.Option) (*gotwilio.SmsResponse, *gotwilio.Exception, error)
}

func main() {
	if err := Run(); err != nil {
		log.Fatal(err)
	}
}

func Run() error {
	log.Info("Starting up...")

	log.Info("Connecting to database...")
	db, err := gorm.Open(sqlite.Open(databaseURL), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("Failed to connect to database: %v\n", err)
	}

	log.Info("Migrating database...")
	if err := db.AutoMigrate(&Reminder{}); err != nil {
		return fmt.Errorf("Failed to migrate reminders table: %v\n", err)
	}
	log.Info("Migrated database")

	log.Info("Initializing SMS client...")
	twilio := gotwilio.NewTwilioClient(twilioAccountSID, twilioAuthToken)
	log.Info("SMS client initialized")

	log.Info("Fetching reminders from database...")
	reminders := []Reminder{}
	if err := db.Find(&reminders).Error; err != nil {
		return fmt.Errorf("Failed to fetch data: %v\n", err)
	}
	log.Info("Fetched reminders from database")

	for tick := range time.Tick(time.Second * 1) {
		for _, reminder := range reminders {
			go func(reminder Reminder) {
				if reminder.MatchesDayAndTime(tick) {
					log.Info("Reminder triggered: ", reminder)
					reminder.SendMessage(twilio, twilioFromNumber, twilioToNumber)
				}
			}(reminder)
		}
	}

	return nil
}
