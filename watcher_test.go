package main

import (
	"context"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/sfreiberg/gotwilio"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/mock"
)

func TestMain(m *testing.M) {
	log.SetOutput(ioutil.Discard)
	os.Exit(m.Run())
}

type SenderMock struct {
	mock.Mock
}

func (m *SenderMock) SendSMS(a string, b string, c string, d string, e string, f ...*gotwilio.Option) (*gotwilio.SmsResponse, *gotwilio.Exception, error) {
	m.Called()
	return nil, nil, nil
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
					Frequency: "Daily",
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
			test.client.On("SendSMS").Return(nil, nil, nil)

			go func() {
				time.Sleep(time.Second * 2)
				cancel()
			}()
			watcher.WatchReminders(test.reminders, test.client)
			test.client.AssertCalled(t, "SendSMS")
		})
	}
}
