package statuscake

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/google/go-querystring/query"
)

type Ssl struct {
	Id             string                 `json:"id"                 url:"id,omitempty"`
	Domain         string              `json:"domain"             url:"domain,omitempty"`
	Checkrate      int                 `url:"checkrate,omitempty"`
	ContactGroupsC string              `url:"contact_groups,omitempty"`
	AlertAt        string              `json:"alert_at"           url:"alert_at,omitempty"`
	AlertReminder  bool                `json:"alert_reminder"     url:"alert_expiry,omitempty"`
	AlertExpiry    bool                `json:"alert_expiry"       url:"alert_reminder,omitempty"`
	AlertBroken    bool                `json:"alert_broken"       url:"alert_broken,omitempty"`
	AlertMixed     bool                `json:"alert_mixed"        url:"alert_mixed,omitempty"`
	Paused         bool                `json:"paused"`
	CertScore      string              `json:"cert_score"`
	CipherScore    string              `json:"cipher_score"`
	CertStatus     string              `json:"cert_status"`
	Cipher         string              `json:"cipher"`
	ValidFromUtc   string              `json:"valid_from_utc"`
	ValidUntilUtc  string              `json:"valid_until_utc"`
	MixedContent   []map[string]string `json:"mixed_content"`
	Flags          map[string]bool     `json:"flags"`
	ContactGroups  []int               `json:"contact_groups"`
	LastReminder   int                 `json:"last_reminder"`
	LastUpdatedUtc string              `json:"last_updated_utc"`
}

type sslUpdateResponse struct {
	Success bool   `json:"Success"`
	Message string `json:"Message"`
}

type Ssls interface {
	All() ([]*Ssl, error)
	Detail(string) (*Ssl, error)
	Update(*Ssl) (*Ssl, error)
	Delete(ID string) error
}

func consolidateSsl(s *Ssl) {
	(*s).ContactGroupsC = strings.Trim(strings.Join(strings.Fields(fmt.Sprint((*s).ContactGroups)), ","), "[]")
}

func findSsl(responses []*Ssl, id string) (*Ssl, error) {
	var response *Ssl
	for _, elem := range responses {
		if (*elem).Id == id {
			return elem, nil
		}
	}
	return response, fmt.Errorf("%s Not found", id)
}

type ssls struct {
	client apiClient
}

func newSsls(c apiClient) Ssls {
	return &ssls{
		client: c,
	}
}

func (tt *ssls) All() ([]*Ssl, error) {
	raw_response, err := tt.client.get("/SSL", nil)
	if err != nil {
		return nil, fmt.Errorf("Error getting StatusCake Ssl: %s", err.Error())
	}
	var getResponse []*Ssl
	err = json.NewDecoder(raw_response.Body).Decode(&getResponse)
	if err != nil {
		return nil, err
	}

	for ssl := range getResponse {
		consolidateSsl(getResponse[ssl])
	}

	return getResponse, err
}

func (tt *ssls) Detail(Id string) (*Ssl, error) {
	responses, err := tt.All()
	if err != nil {
		return nil, err
	}
	mySsl, errF := findSsl(responses, Id)
	if errF != nil {
		return nil, errF
	}
	return mySsl, nil
}

func (tt *ssls) Update(s *Ssl) (*Ssl, error) {
	v, _ := query.Values(*s)
	raw_response, err := tt.client.put("/SSL/Update", v)
	if err != nil {
		return nil, fmt.Errorf("Error creating StatusCake Ssl: %s", err.Error())
	}
	var updateResponse sslUpdateResponse
	err = json.NewDecoder(raw_response.Body).Decode(&updateResponse)
	if err != nil {
		return nil, err
	}

	if !updateResponse.Success {
		return nil, fmt.Errorf("%s", updateResponse.Message)
	}
	return tt.Detail(updateResponse.Message)
}

func (tt *ssls) Delete(id string) error {
	_, err := tt.client.delete("/SSL/Update", url.Values{"id": {fmt.Sprint(id)}})
	if err != nil {
		return err
	}

	return nil
}
