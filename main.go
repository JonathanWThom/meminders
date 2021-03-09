package main

import (
	"context"
	"fmt"
	"net/url"
	"os"

	"github.com/joho/godotenv"
	"github.com/kevinburke/twilio-go"
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

type Sender interface {
	SendMessage(string, string, string, []*url.URL) (*twilio.Message, error)
}

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
}

func main() {
	if err := Run(); err != nil {
		log.Fatal(err)
	}
}

func Run() error {
	log.Info("Starting up...")

	log.Info("Parsing environment...")
	env := os.Getenv("MEMINDERS_ENV")
	if "" == env {
		env = "development"
	}
	if err := godotenv.Load(".env." + env); err != nil {
		log.Fatal("Error loading environment file for env: ", env)
	}

	databaseURL = getenv("DATABASE_URL")
	twilioAccountSID = getenv("TWILIO_ACCOUNT_SID")
	twilioAuthToken = getenv("TWILIO_AUTH_TOKEN")
	twilioFromNumber = getenv("TWILIO_FROM_NUMBER")
	twilioToNumber = getenv("TWILIO_TO_NUMBER")
	log.Info("Environment parsed...")

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

func setUpSMSClient() *twilio.MessageService {
	log.Info("Initializing SMS client...")
	client := twilio.NewClient(twilioAccountSID, twilioAuthToken, nil)
	log.Info("SMS client initialized")

	return client.Messages
}
