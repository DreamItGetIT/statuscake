package statuscake

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakeHTTPClient struct {
	requests []*http.Request
}

func (c *fakeHTTPClient) Do(r *http.Request) (*http.Response, error) {
	c.requests = append(c.requests, r)
	return nil, nil
}

func TestClient(t *testing.T) {
	assert := assert.New(t)

	c := New("random-user", "my-pass")

	assert.Equal("random-user", c.username)
	assert.Equal("my-pass", c.apiKey)
}

func TestClient_newRequest(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	c := New("random-user", "my-pass")
	r, err := c.newRequest("GET", "/hello", nil)

	require.Nil(err)
	assert.Equal("GET", r.Method)
	assert.Equal("https://www.statuscake.com/API/hello", r.URL.String())
	assert.Equal("random-user", r.Header.Get("Username"))
	assert.Equal("my-pass", r.Header.Get("API"))
}

func TestClient_doRequest(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	c := New("random-user", "my-pass")
	hc := &fakeHTTPClient{}
	c.c = hc

	req, err := http.NewRequest("GET", "http://example.com/test", nil)
	require.Nil(err)

	_, err = c.doRequest(req)
	require.Nil(err)

	assert.Len(hc.requests, 1)
	assert.Equal("http://example.com/test", hc.requests[0].URL.String())
}

func TestClient_get(t *testing.T) {
	assert := assert.New(t)

	c := New("random-user", "my-pass")
	hc := &fakeHTTPClient{}
	c.c = hc

	c.get("/hello")
	assert.Len(hc.requests, 1)
	assert.Equal("GET", hc.requests[0].Method)
	assert.Equal("https://www.statuscake.com/API/hello", hc.requests[0].URL.String())
}

func TestClient_Tests(t *testing.T) {
	assert := assert.New(t)

	c := New("foo", "bar")
	expected := &tests{
		client: c,
	}

	assert.Equal(expected, c.Tests())
}
