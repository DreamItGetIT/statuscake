package statuscake

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJSONNumberStringMarshal(t *testing.T) {
	assert := assert.New(t)

	b, err := json.Marshal(jsonNumberString("123"))

	assert.Equal(b, []byte(`"123"`))
	assert.NoError(err)
}

func TestJSONNumberStringUnmarshal(t *testing.T) {
	assert := assert.New(t)

	type testCase struct {
		Input         []byte
		ExpectedValue jsonNumberString
		ExpectedError bool
		ActualValue   jsonNumberString
		ActualError   error
	}
	cases := []testCase{
		{
			Input:         []byte(`"123"`),
			ExpectedValue: jsonNumberString("123"),
		},
		{
			Input:         []byte(`""`),
			ExpectedValue: jsonNumberString(""),
		},
		{
			Input:         []byte(`123`),
			ExpectedValue: jsonNumberString("123"),
		},
		{
			Input:         []byte(`null`),
			ExpectedValue: jsonNumberString(""),
		},
		{
			Input:         []byte(`123.456`),
			ExpectedError: true,
		},
		{
			Input:         []byte(`["123"]`),
			ExpectedError: true,
		},
	}
	for i, c := range cases {
		cases[i].ActualError = json.Unmarshal(c.Input, &cases[i].ActualValue)
	}

	for _, c := range cases {
		if c.ExpectedError {
			assert.Error(c.ActualError)
			continue
		}
		assert.Equal(c.ExpectedValue, c.ActualValue)
		assert.NoError(c.ActualError)
	}
}

func TestTruncate(t *testing.T) {
	assert := assert.New(t)

	assert.Equal([]byte("foo"), truncate([]byte("foo"), 10))
	assert.Equal([]byte("foo"), truncate([]byte("foo"), 3))
	assert.Equal([]byte("foo..."), truncate([]byte("foobarbaz"), 6))
	assert.Equal([]byte("foo..."), truncate([]byte("foobarbaz"), 1))
}
