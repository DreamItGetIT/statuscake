package statuscake

import (
	"fmt"
	"strings"
)

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
