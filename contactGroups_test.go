package statuscake

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/url"
)

func TestContactGroup_All(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	c := &fakeAPIClient{
		fixture: "contactGroupListAllOk.json",
	}
	tt := NewContactGroups(c)
	contactGroups, err := tt.All()
	require.Nil(err)

	assert.Equal("/ContactGroups", c.sentRequestPath)
	assert.Equal("GET", c.sentRequestMethod)
	assert.Nil(c.sentRequestValues)
	assert.Len(contactGroups, 3)
	expectedContactGroup := &ContactGroup{
		GroupName:      "group name",
		Emails:         []string{},
		Mobiles:        "",
		Boxcar:         "",
		Pushover:       "",
		ContactID:      12345,
		DesktopAlert:   "",
		PingURL:        "",
	}
	assert.Equal(expectedContactGroup, contactGroups[0])
	expectedContactGroup.Emails=[]string{"aaaaaaa"}
	expectedContactGroup.PingURL="http"
	expectedContactGroup.ContactID=123456
	assert.Equal(expectedContactGroup, contactGroups[1])
	expectedContactGroup.Emails=[]string{"aaaaaaa","bbbbbbb"}
	expectedContactGroup.ContactID=1234567
	assert.Equal(expectedContactGroup, contactGroups[2])
}

func TestContactGroup_Detail(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	c := &fakeAPIClient{
		fixture: "contactGroupListAllOk.json",
	}
	tt := NewContactGroups(c)

	contactGroup, err := tt.Detail(123456)
	require.Nil(err)
	assert.Equal("/ContactGroups", c.sentRequestPath)
	assert.Equal("GET", c.sentRequestMethod)
	assert.Nil(c.sentRequestValues)
	
	expectedContactGroup := &ContactGroup{
		GroupName:      "group name",
		Emails:         []string{"aaaaaaa"},
		Mobiles:        "",
		Boxcar:         "",
		Pushover:       "",
		ContactID:      123456,
		DesktopAlert:   "",
		PingURL:        "http",
	}
	assert.Equal(expectedContactGroup, contactGroup)
}

func TestContactGroup_Create(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	c := &fakeAPIClient{
		fixture: "contactGroupCreateOk.json",
	}
	tt := NewContactGroups(c)
	contactGroup := &ContactGroup{
		GroupName:      "group name",
		Emails:         []string{"aaaaaa","bbbbbb"},
		PingURL:        "http",
	}
	res, err := tt.Create(contactGroup)
	require.Nil(err)
	assert.Equal("/ContactGroups/Update", c.sentRequestPath)
	assert.Equal("PUT", c.sentRequestMethod)
	assert.Equal(c.sentRequestValues,url.Values(url.Values{"GroupName":[]string{"group name"},"Email":[]string{"aaaaaa,bbbbbb"},"Emails":[]string{"aaaaaa", "bbbbbb"},"PingURL":[]string{"http"},}))
	contactGroup.ContactID=157273
	assert.Equal(contactGroup, res)
}

func TestContactGroup_Update(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	c := &fakeAPIClient{
		fixture: "contactGroupUpdateOk.json",
	}
	tt := NewContactGroups(c)

	contactGroup := &ContactGroup{
		GroupName:      "group name",
		Emails:         []string{"aaaaaa","bbbbbb"},
		PingURL:        "http",
		ContactID:      12345,
	}
	
	res, err := tt.Update(contactGroup)
	require.Nil(err)
	assert.Equal(contactGroup, res)
	assert.Equal("/ContactGroups/Update", c.sentRequestPath)
	assert.Equal("PUT", c.sentRequestMethod)
	assert.Equal(c.sentRequestValues,url.Values(url.Values{"GroupName":[]string{"group name"},"Email":[]string{"aaaaaa,bbbbbb"},"Emails":[]string{"aaaaaa", "bbbbbb"},"PingURL":[]string{"http"},"ContactID":[]string{"12345"},}))
}

func TestContactGroup_Delete(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	c := &fakeAPIClient{
		fixture: "contactGroupDeleteOk.json",
	}
	tt := NewContactGroups(c)

	err := tt.Delete(12345)
	require.Nil(err)
	assert.Equal("/ContactGroups/Update", c.sentRequestPath)
	assert.Equal("DELETE", c.sentRequestMethod)
	assert.Equal(c.sentRequestValues,url.Values(url.Values{"ContactID":[]string{"12345"},}))
}
