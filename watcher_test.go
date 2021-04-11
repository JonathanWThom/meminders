package main

import (
	"context"
	"net/url"
	"testing"
	"time"

	"github.com/jonathanwthom/meminders/mocks"
	"github.com/kevinburke/twilio-go"
)

func TestWatcherWatchReminders(t *testing.T) {
	tests := []struct {
		reminders   []Reminder
		client      Client
		expected    bool
		description string
		calls       *mocks.Caller
		messages    *mocks.Sender
	}{
		{
			reminders: []Reminder{
				{
					Frequency: Daily,
					Hour:      time.Now().Hour(),
					Minute:    time.Now().Minute(),
					Second:    time.Now().Second() + 1,
					Zone:      time.Now().Location().String(),
				},
			},
			client:      &CommsClient{},
			expected:    true,
			description: "Calls sender.SendSMS when reminder is due",
			calls:       &mocks.Caller{},
			messages:    &mocks.Sender{},
		},
		{
			reminders: []Reminder{
				{
					Call:      true,
					Frequency: Daily,
					Hour:      time.Now().Hour(),
					Minute:    time.Now().Minute(),
					Second:    time.Now().Second() + 3,
					Zone:      time.Now().Location().String(),
				},
			},
			client:      &CommsClient{},
			expected:    true,
			description: "Calls caller.Create when reminder is due",
			calls:       &mocks.Caller{},
			messages:    &mocks.Sender{},
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			watcher := Watcher{
				ctx: ctx,
			}
			client := &CommsClient{messages: test.messages, calls: test.calls}
			test.messages.On(
				"SendMessage",
				twilioFromNumber,
				twilioToNumber,
				"",
				[]*url.URL(nil),
			).Return(&twilio.Message{}, nil)
			test.calls.On(
				"Create",
				context.Background(),
				url.Values{"From": []string{twilioFromNumber},
					"To":    []string{twilioToNumber},
					"Twiml": []string{"<Response><Say voice=\"Polly.Joanna\"></Say></Response>"}},
			).Return(&twilio.Call{}, nil)

			go func() {
				time.Sleep(time.Second * 2)
				cancel()
			}()
			watcher.reminders = test.reminders
			watcher.WatchReminders(client)

			// TODO: Could maybe move this to reminders spec
			if test.reminders[0].Call {
				test.calls.AssertCalled(
					t,
					"Create",
					context.Background(),
					url.Values{"From": []string{twilioFromNumber},
						"To":    []string{twilioToNumber},
						"Twiml": []string{"<Response><Say voice=\"Polly.Joanna\"></Say></Response>"}},
				)
			} else {
				test.messages.AssertCalled(
					t,
					"SendMessage",
					twilioFromNumber,
					twilioToNumber,
					"",
					[]*url.URL(nil),
				)
			}
		})
	}
}
