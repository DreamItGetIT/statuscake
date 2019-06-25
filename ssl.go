package statuscake

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"strconv"
	
	"github.com/google/go-querystring/query"
)

type Ssl struct {
	Id             string              `json:"id"                 url:"id,omitempty"`
	Domain         string              `json:"domain"             url:"domain,omitempty"`
	Checkrate      int                 `json:"checkrate"          url:"checkrate,omitempty"`
	ContactGroupsC string              `                          url:"contact_groups,omitempty"`
	AlertAt        string              `json:"alert_at"           url:"alert_at,omitempty"`
	AlertReminder  bool                `json:"alert_reminder"     url:"alert_expiry,omitempty"`
	AlertExpiry    bool                `json:"alert_expiry"       url:"alert_reminder,omitempty"`
	AlertBroken    bool                `json:"alert_broken"       url:"alert_broken,omitempty"`
	AlertMixed     bool                `json:"alert_mixed"        url:"alert_mixed,omitempty"`
	Paused         bool                `json:"paused"`
	IssuerCn       string              `json:"issuer_cn"`
	CertScore      string              `json:"cert_score"`
	CipherScore    string              `json:"cipher_score"`
	CertStatus     string              `json:"cert_status"`
	Cipher         string              `json:"cipher"`
	ValidFromUtc   string              `json:"valid_from_utc"`
	ValidUntilUtc  string              `json:"valid_until_utc"`
	MixedContent   []map[string]string `json:"mixed_content"`
	Flags          map[string]bool     `json:"flags"`
	ContactGroups  []string            `json:"contact_groups"`
	LastReminder   int                 `json:"last_reminder"`
	LastUpdatedUtc string              `json:"last_updated_utc"`
}

// ParialTest represent the a ssl test creation or modification
type PartialSsl struct {
	Id             int
	Domain         string
	Checkrate      string
	ContactGroupsC string
	AlertAt        string
	AlertExpiry    bool
	AlertReminder  bool
	AlertBroken    bool
	AlertMixed     bool
}

type createSsl struct {
	Id             int    `url:"id,omitempty"`
	Domain         string `url:"domain"         json:"domain"`
	Checkrate      string `url:"checkrate"      json:"checkrate"`
	ContactGroupsC string `url:"contact_groups" json:"contact_groups"`
	AlertAt        string `url:"alert_at"       json:"alert_at"`
	AlertExpiry    bool   `url:"alert_expiry"   json:"alert_expiry"`
	AlertReminder  bool   `url:"alert_reminder" json:"alert_reminder"`
	AlertBroken    bool   `url:"alert_broken"   json:"alert_broken"`
	AlertMixed     bool   `url:"alert_mixed"    json:"alert_mixed"`
}

type updateSsl struct {
	Id             int    `url:"id"`
	Domain         string `url:"domain"         json:"domain"`
	Checkrate      string `url:"checkrate"      json:"checkrate"`
	ContactGroupsC string `url:"contact_groups" json:"contact_groups"`
	AlertAt        string `url:"alert_at"       json:"alert_at"`
	AlertExpiry    bool   `url:"alert_expiry"   json:"alert_expiry"`
	AlertReminder  bool   `url:"alert_reminder" json:"alert_reminder"`
	AlertBroken    bool   `url:"alert_broken"   json:"alert_broken"`
	AlertMixed     bool   `url:"alert_mixed"    json:"alert_mixed"`
}


type sslUpdateResponse struct {
	Success bool   `json:"Success"`
	Message interface{} `json:"Message"`
}

type sslCreateResponse struct {
	Success bool   `json:"Success"`
	Message interface{} `json:"Message"`
	Input createSsl `json:"Input"`
}

type Ssls interface {
	All() ([]*Ssl, error)
	completeSsl(*PartialSsl) (*Ssl, error)
	Detail(string) (*Ssl, error)
	Update(*PartialSsl) (*Ssl, error)
	UpdatePartial(*PartialSsl) (*PartialSsl, error)
	Delete(ID string) error
	CreatePartial(*PartialSsl) (*PartialSsl, error)
	Create(*PartialSsl) (*Ssl, error)
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

func (tt *ssls) completeSsl(s *PartialSsl) (*Ssl, error) {
	full, err := tt.Detail(strconv.Itoa((*s).Id))
	if err != nil {
		return nil, err
	}
	(*full).ContactGroups = strings.Split((*s).ContactGroupsC,",")
	return full, nil
}

func Partial(s *Ssl) (*PartialSsl,error) {
	if s==nil {
		return nil,fmt.Errorf("s is nil")
	} else {
		id,err:=strconv.Atoi(s.Id)
		if(err!=nil){
			return nil,err
		}
		return &PartialSsl{
			Id: id,
			Domain: s.Domain,
			Checkrate: strconv.Itoa(s.Checkrate),
			ContactGroupsC: s.ContactGroupsC,
			AlertReminder: s.AlertReminder,
			AlertExpiry: s.AlertExpiry,
			AlertBroken: s.AlertBroken,
			AlertMixed: s.AlertMixed,
			AlertAt: s.AlertAt,
		},nil
	}
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

func (tt *ssls) Update(s *PartialSsl) (*Ssl, error) {
	var err error
	s, err = tt.UpdatePartial(s)
	if err!= nil {
		return nil, err
	}
	return tt.completeSsl(s)
}

func (tt *ssls) UpdatePartial(s *PartialSsl) (*PartialSsl, error) {

	if((*s).Id == 0){
		return tt.CreatePartial(s)
	}
	var v url.Values

	v, _ = query.Values(updateSsl(*s))
	
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
		return nil, fmt.Errorf("%s", updateResponse.Message.(string))
	}


	return s, nil
}

func (tt *ssls) Delete(id string) error {
	_, err := tt.client.delete("/SSL/Update", url.Values{"id": {fmt.Sprint(id)}})
	if err != nil {
		return err
	}

	return nil
}

func (tt *ssls) Create(s *PartialSsl) (*Ssl, error) {
	var err error
	s, err = tt.CreatePartial(s)
	if err!= nil {
		return nil, err
	}
	return tt.completeSsl(s)
}

func (tt *ssls) CreatePartial(s *PartialSsl) (*PartialSsl, error) {
	(*s).Id=0
	var v url.Values
	v, _ = query.Values(createSsl(*s))
	
	raw_response, err := tt.client.put("/SSL/Update", v)
	if err != nil {
		return nil, fmt.Errorf("Error creating StatusCake Ssl: %s", err.Error())
	}

	var createResponse sslCreateResponse
	err = json.NewDecoder(raw_response.Body).Decode(&createResponse)
	if err != nil {
		return nil, err
	}

	if !createResponse.Success {
		return nil, fmt.Errorf("%s", createResponse.Message.(string))
	}
	*s = PartialSsl(createResponse.Input)
	(*s).Id = int(createResponse.Message.(float64))
	
	return s,nil
}

