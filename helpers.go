package main

import "os"

func getenv(name string) string {
	v := os.Getenv(name)
	if v == "" {
		panic("Missing required environment variable: " + name)
	}

	return v
}
