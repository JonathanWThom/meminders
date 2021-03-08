package main

import (
	"fmt"
	"time"
)

type Frequency string

type Reminder struct {
	Message string
	Frequency
	Hour   int
	Minute int
	Second int
	Zone   string
}

func (r *Reminder) MatchesTime(tick time.Time) bool {
	locale, err := time.LoadLocation(r.Zone)
	if err != nil {
		fmt.Println(err)
		return false
	}
	t := tick.In(locale)

	return r.Hour == t.Hour() && r.Minute == t.Minute() && r.Second == t.Second()
}

func (r *Reminder) SendMessage(sender Sender, from string, to string) {
	sender.SendSMS(from, to, r.Message, "", "")
}
