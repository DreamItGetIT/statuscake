package statuscake

import (
  "testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)
func TestSsl_All(t *testing.T) {
  assert := assert.New(t)
	require := require.New(t)

	c := &fakeAPIClient{
		fixture: "ssls_all_ok.json",
	}
	tt := newSsls(c)
	ssls, err := tt.All()
	require.Nil(err)

	assert.Equal("/SSL", c.sentRequestPath)
	assert.Equal("GET", c.sentRequestMethod)
	assert.Nil(c.sentRequestValues)
	assert.Len(ssls, 2)
  mixed := make(map[string]string)
  mixed["type"]="img"
  mixed["src"]="http://example.com/image.gif"
  flags := make(map[string]bool)
  flags["is_extended"] = false
  flags["has_pfs"] = true
  flags["is_broken"] = false
  flags["is_expired"] = false
  flags["is_missing"] = false
  flags["is_revoked"] = false
  flags["is_mixed"] = false
	expectedTest := &Ssl{
    Id: "12345",
    Paused: false,
    Domain: "https://google.com",
    CertScore: "95",
    CipherScore: "100",
    CertStatus: "CERT_OK",
    Cipher: "TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256",
    ValidFromUtc: "2017-10-10 14:06:00",
    ValidUntilUtc: "2017-12-29 00:00:00",
    MixedContent: []map[string]string{mixed},
    Flags: flags,
    ContactGroups: []int{12, 13, 14},
    ContactGroupsC: "12,13,14",
    AlertAt: "1,7,30",
    LastReminder: 0,
    AlertReminder: false,
    AlertExpiry: false,
    AlertBroken: false,
    AlertMixed: false,
    LastUpdatedUtc: "2017-10-24 09:02:25",
	}
	assert.Equal(expectedTest, ssls[0])

	expectedTest = &Ssl{
	  Id: "12346",
    Paused: false,
    Domain: "https://google2.com",
    CertScore: "95",
    CipherScore: "100",
    CertStatus: "CERT_OK",
    Cipher: "TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256",
    ValidFromUtc: "2017-10-10 14:06:00",
    ValidUntilUtc: "2017-12-29 00:00:00",
    MixedContent: []map[string]string{mixed},
    Flags: flags,
    ContactGroups: []int{12, 13, 14},
    ContactGroupsC: "12,13,14",
    AlertAt: "1,7,30",
    LastReminder: 0,
    AlertReminder: false,
    AlertExpiry: false,
    AlertBroken: false,
    AlertMixed: false,
    LastUpdatedUtc: "2017-10-24 09:02:25",
	}
	assert.Equal(expectedTest, ssls[1])
}

func TestSsls_Detail_OK(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	c := &fakeAPIClient{
		fixture: "ssls_all_ok.json",
	}
	tt := newSsls(c)

	ssl, err := tt.Detail("12345")
  require.Nil(err)
	assert.Equal("/SSL", c.sentRequestPath)
	assert.Equal("GET", c.sentRequestMethod)
	assert.Nil(c.sentRequestValues)
  assert.Equal(ssl.Domain, "https://google.com")
  assert.Equal(ssl.Id, "12345")
}

func TestSsls_UpdatePartialCreate_OK(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	c := &fakeAPIClient{
		fixture: "ssls_create_ok.json",
	}
	tt := newSsls(c)
  partial := &PartialSsl{
    Domain: "https://example.com",
  }
  expectedRes := &PartialSsl {
    Id: 12345,
    Domain: "https://example.com",
    Checkrate: 86400,
    ContactGroupsC: "1000,2000",
    AlertReminder: false,
    AlertExpiry: false,
    AlertBroken: false,
    AlertAt: "59,60,61",
  }
  res, err := tt.UpdatePartial(partial)
  require.Nil(err)
  assert.Equal(expectedRes, res)
}

func TestSsls_UpdatePartial_OK(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	c := &fakeAPIClient{
		fixture: "ssls_update_ok.json",
	}
	tt := newSsls(c)
  partial := &PartialSsl{
    Id: 12345,
    AlertReminder: true,
  }
  expectedRes := &PartialSsl {
    Id: 12345,
    AlertReminder: true,
  }
  res, err := tt.UpdatePartial(partial)
  require.Nil(err)
  assert.Equal(expectedRes, res)
}
func TestSsl_complete_OK(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	c := &fakeAPIClient{
		fixture: "ssls_all_ok.json",
	}
	tt := newSsls(c)

  partial := &PartialSsl {
    Id: 12345,
    Domain: "https://example.com",
    Checkrate: 86400,
    ContactGroupsC: "1000,2000",
    AlertReminder: false,
    AlertExpiry: false,
    AlertBroken: false,
    AlertAt: "59,60,61",
  }
  full, err := tt.completeSsl(partial)
  require.Nil(err)
  assert.Equal(full.Domain, "https://google.com")
  assert.Equal(full.CipherScore, "100")
  assert.Equal(full.Checkrate, 86400)
  assert.Equal(full.ContactGroups, []int{1000, 2000})
}
