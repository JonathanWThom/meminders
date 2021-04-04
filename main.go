package main

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"time"

	"github.com/die-net/http-tarpit/tarpit"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/kevinburke/twilio-go"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	log "github.com/sirupsen/logrus"
)

var (
	adminPassword,
	adminUsername,
	databaseURL,
	twilioAccountSID,
	twilioAuthToken,
	twilioFromNumber,
	twilioToNumber string
)

// TODO: Refactor many of these global vars into config if possible

var watcher = Watcher{
	ctx: context.Background(),
}

var db *gorm.DB

var tp *tarpit.Tarpit

type Router interface {
	Run(...string) error
	ServeHTTP(http.ResponseWriter, *http.Request)
}

type Config interface {
	buildRouter() Router
}

type config struct{}

var appConfig Config

type CommsClient struct {
	Messages Sender
	Calls    Caller
}

type Sender interface {
	SendMessage(string, string, string, []*url.URL) (*twilio.Message, error)
}

type Caller interface {
	Create(context.Context, url.Values) (*twilio.Call, error)
}

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
}

func main() {
	appConfig = &config{}
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

	adminPassword = getenv("ADMIN_PASSWORD")
	adminUsername = getenv("ADMIN_USERNAME")
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

	client := setUpCommunicationsClient()

	reminders, err := getReminders(db)
	if err != nil {
		return err
	}
	watcher.reminders = reminders
	go watcher.WatchReminders(client)

	router := appConfig.buildRouter()
	log.Info("Routes built and serving on port :8080")
	err = router.Run()
	if err != nil {
		return err
	}

	return nil
}

func (c *config) buildRouter() Router {
	log.Info("Building routes...")
	router := gin.Default()
	authorized := router.Group("/", gin.BasicAuth(gin.Accounts{
		adminUsername: adminPassword,
	}))
	authorized.POST("/reminders", postReminders)
	router.Any("/robots.txt", robotsHandler)

	// Tarpit
	tp = tarpit.New(
		runtime.NumCPU(),
		"text/html",
		16*time.Second,
		50*time.Millisecond,
		1048576,
		10485760,
	)
	if tp == nil {
		log.Fatal("Unable to build tarpit")
	}
	router.NoRoute(tarpitHandler)

	return router
}

func robotsHandler(ctx *gin.Context) {
	ctx.String(http.StatusOK, "User-agent: * Disallow /")
}

func tarpitHandler(ctx *gin.Context) {
	tp.Handler(ctx.Writer, ctx.Request)
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

func setUpCommunicationsClient() *CommsClient {
	log.Info("Initializing Communications client...")
	client := twilio.NewClient(twilioAccountSID, twilioAuthToken, nil)
	log.Info("Communcations client initialized")

	return &CommsClient{Calls: client.Calls, Messages: client.Messages}
}
