package statuscake

import (
	"fmt"
	"strings"
)

// APIError implements the error interface an it's used when the API response has errors.
type APIError interface {
	APIError() string
}

type httpError struct {
	status     string
	statusCode int
}

func (e *httpError) Error() string {
	return fmt.Sprintf("HTTP error: %d %s", e.statusCode, e.status)
}

// ValidationError is a map where the key is the invalid field and the value is a message describing why the field is invalid.
type ValidationError map[string]string

func (e ValidationError) Error() string {
	var messages []string

	for k, v := range e {
		m := fmt.Sprintf("%s %s", k, v)
		messages = append(messages, m)
	}

	return strings.Join(messages, ", ")
}

type updateError struct {
	Issues map[string]string
}

func (e *updateError) Error() string {
	var messages []string

	for k, v := range e.Issues {
		m := fmt.Sprintf("%s %s", k, v)
		messages = append(messages, m)
	}

	return strings.Join(messages, ", ")
}

// ApiError returns the error specified in the API response
func (e *updateError) ApiError() string {
	return e.Error()
}
