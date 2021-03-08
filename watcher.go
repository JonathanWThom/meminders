package main

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"
)

type Watcher struct {
	ctx context.Context
}

func (w *Watcher) WatchReminders(reminders []Reminder, client Sender) {
ticker:
	for tick := range time.Tick(time.Second * 1) {
		select {
		case <-w.ctx.Done():
			break ticker
		default:
			for _, reminder := range reminders {
				go func(reminder Reminder) {
					if reminder.MatchesDayAndTime(tick) {
						log.Info("Reminder triggered: ", reminder)
						reminder.SendMessage(client, twilioFromNumber, twilioToNumber)
					}
				}(reminder)
			}
		}
	}
}
