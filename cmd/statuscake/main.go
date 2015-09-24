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
		fmt.Printf("* %d %s: %s\n", t.TestID, t.WebsiteName, t.Status)
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
