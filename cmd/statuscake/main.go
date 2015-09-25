package main

import (
	"fmt"
	logpkg "log"
	"os"

	"github.com/DreamItGetIT/statuscake"
)

var log *logpkg.Logger

func init() {
	log = logpkg.New(os.Stderr, "", 0)
}

func colouredStatus(s string) string {
	switch s {
	case "Up":
		return fmt.Sprintf("\033[0;32m%s\033[0m", s)
	case "Down":
		return fmt.Sprintf("\033[0;31m%s\033[0m", s)
	default:
		return s
	}
}

func getEnv(name string) string {
	v := os.Getenv(name)
	if v == "" {
		log.Fatalf("`%s` env variable is required", name)
	}

	return v
}

func listTests(c *statuscake.Client) error {
	tt := c.Tests()
	tests, err := tt.All()
	if err != nil {
		return err
	}

	for _, t := range tests {
		var paused string
		if t.Paused {
			paused = "yes"
		} else {
			paused = "no"
		}

		fmt.Printf("* %d: %s\n", t.TestID, colouredStatus(t.Status))
		fmt.Printf("  WebsiteName: %s\n", t.WebsiteName)
		fmt.Printf("  TestType: %s\n", t.TestType)
		fmt.Printf("  Paused: %s\n", paused)
		fmt.Printf("  ContactGroup: %d\n", t.ContactGroup)
		fmt.Printf("  ContactID: %d\n", t.ContactID)
		fmt.Printf("  Uptime: %f\n", t.Uptime)
	}

	return nil
}

func main() {
	username := getEnv("STATUSCAKE_USERNAME")
	apikey := getEnv("STATUSCAKE_APIKEY")

	c := statuscake.New(username, apikey)
	err := listTests(c)
	if err != nil {
		log.Fatal(err)
	}
}
