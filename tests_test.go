package statuscake

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakeApiClient struct {
	requestedPath string
	fixture       string
}

func (c *fakeApiClient) get(path string) (*http.Response, error) {
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

func TestTests_All(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	c := &fakeApiClient{
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
