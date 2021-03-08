package main

import (
	"fmt"
	"time"
)

type Frequency string

// add some validations
type Reminder struct {
	Message string
	Frequency
	Day       int
	DayOfWeek string
	Hour      int
	Minute    int
	Second    int
	Zone      string
}

func (r *Reminder) MatchesDay(t time.Time) bool {
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

func (r *Reminder) MatchesTime(t time.Time) bool {
	return r.Hour == t.Hour() && r.Minute == t.Minute() && r.Second == t.Second()
}

func (r *Reminder) MatchesDayAndTime(tick time.Time) bool {
	locale, err := time.LoadLocation(r.Zone)
	if err != nil {
		fmt.Println(err)
		return false
	}
	t := tick.In(locale)

	return r.MatchesDay(t) && r.MatchesTime(t)
}

func (r *Reminder) SendMessage(sender Sender, from string, to string) {
	sender.SendSMS(from, to, r.Message, "", "")
}
