package main

import (
	"testing"
	"time"
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
