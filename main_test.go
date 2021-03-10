package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	log "github.com/sirupsen/logrus"
)

func TestMain(m *testing.M) {
	os.Setenv("MEMINDERS_ENV", "test")
	log.SetOutput(ioutil.Discard)
	os.Exit(m.Run())
}

type mockConfig struct{}

func (m *mockConfig) buildRouter() Router {
	return &mockRouter{}
}

type mockRouter struct{}

func (m *mockRouter) Run(...string) error {
	return nil
}
func (m *mockRouter) ServeHTTP(http.ResponseWriter, *http.Request) {}

func TestRun(t *testing.T) {
	originalWatcher := watcher
	defer func() { watcher = originalWatcher }()
	oldAppConfig := appConfig
	defer func() { appConfig = oldAppConfig }()
	appConfig = &mockConfig{}

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

func TestRouter(t *testing.T) {
	oldAppConfig := appConfig
	defer func() { appConfig = oldAppConfig }()
	appConfig = &config{}

	tests := []struct {
		authorization bool
		description   string
		headers       map[string]string
		method        string
		path          string
		statusCode    int
	}{
		{
			authorization: true,
			description:   "Valid auth headers are passed, but no data",
			method:        "POST",
			path:          "/reminders",
			statusCode:    400,
		},
		{
			authorization: false,
			description:   "No auth headers are passed",
			method:        "POST",
			path:          "/reminders",
			statusCode:    401,
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			srv := httptest.NewServer(appConfig.buildRouter())
			defer srv.Close()

			client := &http.Client{}
			url := fmt.Sprintf("%s%s", srv.URL, test.path)
			req, err := http.NewRequest(test.method, url, nil)
			if err != nil {
				t.Fatalf("could not send %v request to path %v: %v", test.method, test.path, err)
			}

			if test.authorization {
				req.SetBasicAuth(os.Getenv("ADMIN_USERNAME"), os.Getenv("ADMIN_PASSWORD"))
				fmt.Println(req.Header)
			}

			res, err := client.Do(req)
			if err != nil {
				t.Fatalf("could not send %v request to path %v: %v", test.method, test.path, err)
			}

			if res.StatusCode != test.statusCode {
				t.Errorf(
					"Expected %v request to path %v to have status %v; got %v",
					test.method,
					test.path,
					test.statusCode,
					res.StatusCode,
				)
			}
		})
	}
}
