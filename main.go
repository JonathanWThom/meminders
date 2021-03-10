package main

import (
	"context"
	"fmt"
	"net/url"
	"os"

	"github.com/gin-gonic/gin"
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

var db *gorm.DB

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
	if env != "development" {
		gin.SetMode(gin.ReleaseMode)
	}
	log.Info("Environment parsed...")

	var err error
	db, err = setUpDB()
	if err != nil {
		return err
	}

	reminders, err := getReminders(db)
	if err != nil {
		return err
	}

	client := setUpSMSClient()

	// extract me to package or function
	log.Info("Building routes...")
	router := gin.Default()
	authorized := router.Group("/", gin.BasicAuth(gin.Accounts{
		"foo": "bar", // change me
	}))
	authorized.POST("/reminders", postReminders)
	port := ":8080"
	// don't actually connect for tests
	go router.Run(port)
	log.Info("Routes built and serving on port: ", port)

	watcher.WatchReminders(reminders, client)

	return nil
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
