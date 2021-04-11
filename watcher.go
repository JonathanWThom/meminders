package main

import (
	"context"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
)

type Watcher struct {
	ctx       context.Context
	reminders []Reminder
}

func (w *Watcher) WatchReminders(client Client) {
ticker:
	for tick := range time.Tick(time.Second * 1) {
		select {
		case <-w.ctx.Done():
			break ticker
		default:
			for _, reminder := range w.reminders {
				go func(reminder Reminder) {
					if reminder.MatchesDayAndTime(tick) {
						fmt.Println(reminder.Call)
						log.Info("Reminder triggered: ", reminder)
						reminder.SendMessage(client, twilioFromNumber, twilioToNumber)
					}
				}(reminder)
			}
		}
	}
}
