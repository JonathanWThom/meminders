package main

import (
	"context"
	"io/ioutil"
	"os"
	"testing"

	log "github.com/sirupsen/logrus"
)

func TestMain(m *testing.M) {
	databaseURL = "file::memory:?cache=shared"
	twilioAccountSID = ""
	twilioAuthToken = ""
	twilioFromNumber = ""
	twilioToNumber = ""

	log.SetOutput(ioutil.Discard)
	os.Exit(m.Run())
}

func TestRun(t *testing.T) {
	originalWatcher := watcher
	defer func() { watcher = originalWatcher }()

	tests := []struct {
		description string
		err         error
	}{
		{
			description: "It runs setup without fatal error",
			err:         nil,
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			watcher = Watcher{
				ctx: ctx,
			}

			go cancel()
			err := Run()
			if err != test.err {
				t.Errorf("Expected Run() error to eq %v, got %v", test.err, err)
			}
		})
	}
}
