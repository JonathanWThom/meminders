package main

import (
	"time"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// add some validations
type Reminder struct {
	gorm.Model

	Message   string
	Frequency string
	Day       int
	DayOfWeek string
	Hour      int
	Minute    int
	Second    int
	Zone      string
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
	_, _, err := sender.SendSMS(from, to, r.Message, "", "")
	if err != nil {
		log.Error("Failed to send reminder: ", err)
		return
	}

	log.Info("Reminder sent: ", *r)
}
