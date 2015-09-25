package statuscake

import (
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTest_Validate(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	test := &Test{
		Timeout:      200,
		Confirmation: 100,
		Public:       200,
		Virus:        200,
		TestType:     "FTP",
		RealBrowser:  100,
		TriggerRate:  100,
		CheckRate:    100000,
		WebsiteName:  "",
		WebsiteURL:   "",
	}

	err := test.Validate()
	require.NotNil(err)

	message := err.Error()
	assert.Contains(message, "WebsiteName is required")
	assert.Contains(message, "WebsiteURL is required")
	assert.Contains(message, "Timeout must be 0 or between 6 and 99")
	assert.Contains(message, "Confirmation must be between 0 and 9")
	assert.Contains(message, "CheckRate must be between 0 and 23999")
	assert.Contains(message, "Public must be 0 or 1")
	assert.Contains(message, "Virus must be 0 or 1")
	assert.Contains(message, "TestType must be HTTP, TCP, or PING")
	assert.Contains(message, "RealBrowser must be 0 or 1")
	assert.Contains(message, "TriggerRate must be between 0 and 59")

	test.Timeout = 10
	test.Confirmation = 2
	test.Public = 1
	test.Virus = 1
	test.TestType = "HTTP"
	test.RealBrowser = 1
	test.TriggerRate = 50
	test.CheckRate = 10
	test.WebsiteName = "Foo"
	test.WebsiteURL = "http://example.com"

	err = test.Validate()
	assert.Nil(err)
}

func TestTest_ToURLValues(t *testing.T) {
	assert := assert.New(t)

	test := &Test{
		TestID:      123,
		Paused:      true,
		WebsiteName: "Foo Bar",
	}

	expected := url.Values{
		"TestID":      {"123"},
		"WebsiteName": {"Foo Bar"},
		"Paused":      {"1"},
	}

	assert.Equal(expected, test.ToURLValues())
}

func TestTests_All(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	c := &fakeAPIClient{
		fixture: "tests_ok.json",
	}
	tt := newTests(c)
	tests, err := tt.All()
	require.Nil(err)

	assert.Equal("/Tests", c.requestedPath)
	assert.Len(tests, 2)

	expectedTest := &Test{
		TestID:      100,
		Paused:      false,
		TestType:    "HTTP",
		WebsiteName: "www 1",
		ContactID:   1,
		Status:      "Up",
		Uptime:      100,
	}
	assert.Equal(expectedTest, tests[0])

	expectedTest = &Test{
		TestID:      101,
		Paused:      true,
		TestType:    "HTTP",
		WebsiteName: "www 2",
		ContactID:   2,
		Status:      "Down",
		Uptime:      0,
	}
	assert.Equal(expectedTest, tests[1])
}

type fakeAPIClient struct {
	requestedPath string
	fixture       string
}

func (c *fakeAPIClient) get(path string) (*http.Response, error) {
	c.requestedPath = path
	p := filepath.Join("fixtures", c.fixture)
	f, err := os.Open(p)
	if err != nil {
		log.Fatal(err)
	}

	resp := &http.Response{
		Body: f,
	}

	return resp, nil
}
