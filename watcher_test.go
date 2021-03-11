package main

import (
	"context"
	"net/url"
	"testing"
	"time"

	"github.com/kevinburke/twilio-go"
	"github.com/stretchr/testify/mock"
)

type SenderMock struct {
	mock.Mock
}

func (m *SenderMock) SendMessage(string, string, string, []*url.URL) (*twilio.Message, error) {
	m.Called()
	return &twilio.Message{}, nil
}

func TestWatcherWatchReminders(t *testing.T) {
	tests := []struct {
		reminders   []Reminder
		client      *SenderMock
		expected    bool
		description string
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
			client:      &SenderMock{},
			expected:    true,
			description: "Calls sender.SendSMS when reminder is due",
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			watcher := Watcher{
				ctx: ctx,
			}
			test.client.On("SendMessage").Return(&twilio.Message{}, nil)

			go func() {
				time.Sleep(time.Second * 2)
				cancel()
			}()
			watcher.reminders = test.reminders
			watcher.WatchReminders(test.client)
			test.client.AssertCalled(t, "SendMessage")
		})
	}
}
