package main

import (
	"fmt"
	logpkg "log"
	"os"
	"strconv"

	"github.com/DreamItGetIT/statuscake"
)

var log *logpkg.Logger

type command func(*statuscake.Client, ...string) error

var commands map[string]command

func init() {
	log = logpkg.New(os.Stderr, "", 0)
	commands = map[string]command{
		"list":   cmdList,
		"delete": cmdDelete,
	}
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

func cmdList(c *statuscake.Client, args ...string) error {
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

func cmdDelete(c *statuscake.Client, args ...string) error {
	if len(args) != 1 {
		return fmt.Errorf("command `delete` requires a single argument `TestID`")
	}

	id, err := strconv.Atoi(args[0])
	if err != nil {
		return err
	}

	return c.Tests().Delete(id)
}

func usage() {
	fmt.Printf("Usage:\n")
	fmt.Printf("  %s COMMAND\n", os.Args[0])
	fmt.Printf("Available commands:\n")
	for k, _ := range commands {
		fmt.Printf("  %+v\n", k)
	}
}

func main() {
	username := getEnv("STATUSCAKE_USERNAME")
	apikey := getEnv("STATUSCAKE_APIKEY")

	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}

	var err error

	c := statuscake.New(username, apikey)
	if cmd, ok := commands[os.Args[1]]; ok {
		err = cmd(c, os.Args[2:]...)
	} else {
		err = fmt.Errorf("Unknown command `%s`", os.Args[1])
	}

	if err != nil {
		log.Fatalf("Error running command `%s`: %s", os.Args[1], err.Error())
	}
}
