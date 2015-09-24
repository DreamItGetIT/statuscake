package statuscake

import "encoding/json"

type Test struct {
	TestID      int
	Paused      bool
	TestType    string
	WebsiteName string
	ContactID   int
	Status      string
	Uptime      int
}

type Tests interface {
	All() ([]*Test, error)
}

type tests struct {
	client apiClient
}

func newTests(c apiClient) Tests {
	return &tests{
		client: c,
	}
}

func (t *tests) All() ([]*Test, error) {
	resp, err := t.client.get("/Tests")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var tests []*Test
	err = json.NewDecoder(resp.Body).Decode(&tests)

	return tests, err
}
