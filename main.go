package main

import (
	"context"
	"fmt"
	"os"

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

var watcher = Watcher{
	ctx: context.Background(),
}

func init() {
	env := os.Getenv("MEMINDERS_ENV")
	if "" == env {
		env = "development"
	}
	godotenv.Load(".env." + env)
	godotenv.Load()

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

	db, err := setUpDB()
	if err != nil {
		return err
	}

	reminders, err := getReminders(db)
	if err != nil {
		return err
	}

	client := setUpSMSClient()

	watcher.WatchReminders(reminders, client)

	return nil
}

func getReminders(db *gorm.DB) ([]Reminder, error) {
	log.Info("Fetching reminders from database...")
	reminders := []Reminder{}
	if err := db.Find(&reminders).Error; err != nil {
		return reminders, fmt.Errorf("Failed to fetch data: %v\n", err)
	}
	log.Info("Fetched reminders from database")

	return reminders, nil
}

func setUpDB() (*gorm.DB, error) {
	log.Info("Connecting to database " + databaseURL + "...")
	db, err := gorm.Open(sqlite.Open(databaseURL), &gorm.Config{})
	if err != nil {
		return db, fmt.Errorf("Failed to connect to database: %v\n", err)
	}

	log.Info("Migrating database...")
	if err := db.AutoMigrate(&Reminder{}); err != nil {
		return db, fmt.Errorf("Failed to migrate reminders table: %v\n", err)
	}
	log.Info("Migrated database")

	return db, nil
}

func setUpSMSClient() *gotwilio.Twilio {
	log.Info("Initializing SMS client...")
	twilio := gotwilio.NewTwilioClient(twilioAccountSID, twilioAuthToken)
	log.Info("SMS client initialized")

	return twilio
}
