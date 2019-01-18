package statuscake

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// APIError implements the error interface an it's used when the API response has errors.
type APIError interface {
	APIError() string
}

// HTTPError represents a response returned from the StatusCake API which
// contained an erronous status code.
type HTTPError struct {
	Status     string
	StatusCode int
	Message    string
	ErrorNo    int
}

// HTTPErrorResponse represents the body content of an error response.
type HTTPErrorResponse struct {
	ErrorMessage string `json:"Error,omitempty"`
	ErrorNo      int    `json:"ErrNo,omitempty"`
}

// NewHTTPError returns an HTTPError object from the given http.Response object.
func NewHTTPError(r *http.Response) *HTTPError {
	httpError := HTTPError{}

	if r == nil {
		return &httpError
	}

	// Set the status
	httpError.Status = r.Status
	httpError.StatusCode = r.StatusCode

	// Attempt to read the body to find the error message.  Return if no body exists.
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return &httpError
	}
	var body HTTPErrorResponse
	if err = json.Unmarshal(b, &body); err != nil {
		return &httpError
	}
	httpError.Message = body.ErrorMessage
	httpError.ErrorNo = body.ErrorNo

	return &httpError
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("HTTP error: %d - %s %s", e.StatusCode, e.Status, e.Message)
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
	Issues  interface{}
	Message string
}

func (e *updateError) Error() string {
	var messages []string

	messages = append(messages, e.Message)

	if issues, ok := e.Issues.(map[string]interface{}); ok {
		for k, v := range issues {
			m := fmt.Sprintf("%s %s", k, v)
			messages = append(messages, m)
		}
	} else if issues, ok := e.Issues.([]interface{}); ok {
		for _, v := range issues {
			m := fmt.Sprint(v)
			messages = append(messages, m)
		}
	} else if issue, ok := e.Issues.(interface{}); ok {
		m := fmt.Sprint(issue)
		messages = append(messages, m)
	}

	return strings.Join(messages, ", ")
}

// APIError returns the error specified in the API response
func (e *updateError) APIError() string {
	return e.Error()
}

type deleteError struct {
	Message string
}

func (e *deleteError) Error() string {
	return e.Message
}

// AuthenticationError implements the error interface and it's returned
// when API responses have authentication errors
type AuthenticationError struct {
	errNo   int
	message string
}

func (e *AuthenticationError) Error() string {
	return fmt.Sprintf("%d, %s", e.errNo, e.message)
}
