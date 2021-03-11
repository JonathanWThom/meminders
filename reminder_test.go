package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func TestReminderMatchesDayAndTime(t *testing.T) {
	now := time.Now()
	year := now.Year()
	month := now.Month()
	location := now.Location()
	zone := location.String()

	tests := []struct {
		description string
		expected    bool
		reminder    Reminder
		time        time.Time
	}{
		{
			description: "Reminder is Daily and time matches",
			expected:    true,
			reminder:    Reminder{Frequency: "Daily", Hour: 1, Minute: 1, Second: 1, Zone: zone},
			time:        time.Date(year, month, now.Day(), 1, 1, 1, 0, location),
		},
		{
			description: "Reminder is Daily and time does not match",
			expected:    false,
			reminder:    Reminder{Frequency: "Daily", Hour: 1, Minute: 1, Second: 1, Zone: zone},
			time:        time.Date(year, month, now.Day(), 1, 2, 1, 0, location),
		},
		{
			description: "Reminder is Weekly and time and day of week both match",
			expected:    true,
			reminder:    Reminder{Frequency: "Weekly", DayOfWeek: now.Weekday().String(), Hour: 1, Minute: 1, Second: 1, Zone: zone},
			time:        time.Date(year, month, now.Day(), 1, 1, 1, 0, location),
		},
		{
			description: "Reminder is Weekly and time matches but day of week does not",
			expected:    false,
			reminder:    Reminder{Frequency: "Weekly", DayOfWeek: now.Weekday().String(), Hour: 1, Minute: 1, Second: 1, Zone: zone},
			time:        time.Date(year, month, now.Day()+1, 1, 1, 1, 0, location),
		},
		{
			description: "Reminder is Monthly and time and day of month both match",
			expected:    true,
			reminder:    Reminder{Frequency: "Monthly", Day: now.Day(), Hour: 1, Minute: 1, Second: 1, Zone: zone},
			time:        time.Date(year, month, now.Day(), 1, 1, 1, 0, location),
		},
		{
			description: "Reminder is Weekly and time matches but day of month does not",
			expected:    false,
			reminder:    Reminder{Frequency: "Monthly", Day: now.Day(), Hour: 1, Minute: 1, Second: 1, Zone: zone},
			time:        time.Date(year, month, now.Day()+1, 1, 1, 1, 0, location),
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			actual := test.reminder.MatchesDayAndTime(test.time)

			if actual != test.expected {
				t.Errorf(
					"Reminder.MatchesDayAndTime(%v) returned %v, expected %v",
					test.time,
					actual,
					test.expected,
				)
			}
		})
	}
}

func TestPostReminders(t *testing.T) {
	tests := []struct {
		description    string
		params         map[string]interface{}
		statusCode     int
		remindersCount int
	}{
		{
			description: "It returns 201 and creates a reminder when valid parameters are passed",
			params: map[string]interface{}{
				"message":   "test message",
				"frequency": "Daily",
				"hour":      1,
				"minute":    2,
				"zone":      "America/Los_Angeles",
			},
			statusCode:     201,
			remindersCount: 1,
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			preReminders := []Reminder{}
			db.Find(&preReminders)

			payload, err := json.Marshal(test.params)
			if err != nil {
				t.Fatalf("Failed to build payload for params: %v", test.params)
			}

			req, err := http.NewRequest("POST", "localhost:8080/reminders", bytes.NewBuffer(payload))
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			rec := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rec)
			c.Request = req

			postReminders(c)

			res := rec.Result()
			defer res.Body.Close()

			if res.StatusCode != test.statusCode {
				t.Errorf("postReminders expected status code %v, got %v", test.statusCode, res.StatusCode)
			}

			reminders := []Reminder{}
			db.Find(&reminders)
			remindersAdded := len(reminders) - len(preReminders)
			if remindersAdded != test.remindersCount {
				t.Errorf("postReminders expected %v reminders,  got %v", test.remindersCount, len(reminders))
			}
		})
	}
}
