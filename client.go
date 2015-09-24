package statuscake

import (
	"fmt"
	"io"
	"net/http"
)

const apiBaseURL = "https://www.statuscake.com/API"

type httpClient interface {
	Do(*http.Request) (*http.Response, error)
}

type apiClient interface {
	get(string) (*http.Response, error)
}

type Client struct {
	c        httpClient
	username string
	apiKey   string
}

func New(username string, apiKey string) *Client {
	c := &http.Client{}
	return &Client{
		c:        c,
		username: username,
		apiKey:   apiKey,
	}
}

func (c *Client) newRequest(method string, path string, body io.Reader) (*http.Request, error) {
	url := fmt.Sprintf("%s%s", apiBaseURL, path)
	r, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	r.Header.Set("Username", c.username)
	r.Header.Set("API", c.apiKey)

	return r, nil
}

func (c *Client) doRequest(r *http.Request) (*http.Response, error) {
	return c.c.Do(r)
}

func (c *Client) get(path string) (*http.Response, error) {
	r, err := c.newRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	return c.doRequest(r)
}

func (c *Client) Tests() Tests {
	return newTests(c)
}
