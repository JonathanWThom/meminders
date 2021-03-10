package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// add some validations
// make required?
type Reminder struct {
	gorm.Model

	Message   string `json:"message" binding:"required"`
	Frequency string `json:"frequency" binding:"required"`
	Day       int    `json:"day"`
	DayOfWeek string `json:"day_of_week"`
	Hour      int    `json:"hour" binding:"required"`
	Minute    int    `json:"minute"`
	Second    int    `json:"second"`
	Zone      string `json:"zone" binding:"required"`
}

func postReminders(c *gin.Context) {
	var reminder Reminder
	if err := c.ShouldBindJSON(&reminder); err != nil {
		log.Errorf("Invalid arguments to POST /reminders: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := db.Create(&reminder).Error; err != nil {
		log.Errorf("Failed to create reminder: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create reminder"})
	}

	c.JSON(http.StatusCreated, &reminder)
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

func (r *Reminder) matchesDay(t time.Time) bool {
	if r.Frequency == "Daily" {
		return true
	}

	if r.Frequency == "Weekly" {
		return r.DayOfWeek == t.Weekday().String()
	}

	if r.Frequency == "Monthly" {
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

func (r *Reminder) SendMessage(sender Sender, from string, to string) {
	_, err := sender.SendMessage(from, to, r.Message, nil)
	if err != nil {
		log.Error("Failed to send reminder: ", err)
		return
	}

	log.Info("Reminder sent: ", r)
}
