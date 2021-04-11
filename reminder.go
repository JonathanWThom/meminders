package main

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

const (
	Daily   = "daily"
	Monthly = "monthly"
	Once    = "once"
	Weekly  = "weekly"
)

var frequencies = []string{
	Daily,
	Monthly,
	Once,
	Weekly,
}

type Reminder struct {
	gorm.Model

	Call      bool   `json:"call"`
	Day       int    `json:"day"`
	DayOfWeek string `json:"day_of_week"`
	Frequency string `json:"frequency" binding:"required"`
	Hour      int    `json:"hour" binding:"required"`
	Message   string `json:"message" binding:"required"`
	Minute    int    `json:"minute"`
	Month     string `json:"month"`
	Second    int    `json:"second"`
	Year      int    `json:"year"`
	Zone      string `json:"zone" binding:"required"`
}

func postReminders(c *gin.Context) {
	var reminder Reminder
	if err := c.ShouldBindJSON(&reminder); err != nil {
		log.Errorf("Invalid arguments to POST /reminders: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}

	// Custom Validations
	// TODO: Break these out into a function or use as tag if possible
	// TODO: Validate right params given frequency
	reminder.Frequency = strings.ToLower(reminder.Frequency)
	_, ok := find(frequencies, reminder.Frequency)
	if !ok {
		err := fmt.Errorf("Invalid frequency: %v", reminder.Frequency)
		log.Errorf("Failed to create reminder: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.Create(&reminder).Error; err != nil {
		log.Errorf("Failed to create reminder: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create reminder"})
	}

	c.JSON(http.StatusCreated, &reminder)
	resetWatcher()
}

func resetWatcher() {
	reminders, err := getReminders(db)
	if err != nil {
		log.Error("Failed to refetch reminders: ", err)
		return
	}
	watcher.reminders = reminders
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

func (r *Reminder) matchesSingleDay(t time.Time) bool {
	return r.Day == t.Day() && r.Month == t.Month().String() && r.Year == t.Year()
}

func (r *Reminder) matchesDay(t time.Time) bool {
	if r.Frequency == Once {
		return r.matchesSingleDay(t)
	}

	if r.Frequency == Daily {
		return true
	}

	if r.Frequency == Weekly {
		return r.DayOfWeek == t.Weekday().String()
	}

	if r.Frequency == Monthly {
		return r.Day == t.Day()
	}

	return false
}

func (r *Reminder) matchesTime(t time.Time) bool {
	return r.Hour == t.Hour() && r.Minute == t.Minute() && r.Second == t.Second()
}

func (r *Reminder) MatchesDayAndTime(tick time.Time) bool {
	locale, err := time.LoadLocation(r.Zone)
	if err != nil {
		log.Error("Failed to set local time: ", err)
		return false
	}
	t := tick.In(locale)

	return r.matchesDay(t) && r.matchesTime(t)
}

func (r *Reminder) SendMessage(client Client, from string, to string) {
	fmt.Println(r.Call)
	if r.Call == true {
		data := url.Values{}
		data.Set("From", twilioFromNumber)
		data.Set("To", twilioToNumber)
		msg := fmt.Sprintf(
			"<Response><Say voice=\"Polly.Joanna\">%s</Say></Response>",
			r.Message,
		)
		data.Set("Twiml", msg)
		client.Calls().Create(context.Background(), data)
	} else {
		_, err := client.Messages().SendMessage(from, to, r.Message, nil)
		if err != nil {
			log.Error("Failed to send reminder: ", err)
			return
		}
	}

	log.Info("Reminder sent: ", r)
}
