package response

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"encoding/json"
)

func TestDefaultValues(t *testing.T) {
	ps := NewPingSummary()
	assert.Equal(t, "", ps.Service, "Service should return an empty string as default")
	assert.Equal(t, "", ps.Storage, "Storage should return an empty string as default")
}

func TestServiceOnline(t *testing.T)  {
	ps := NewPingSummary()
	ps.TurnServiceOnline()
	assert.Equal(t, PONG, ps.Service, "As the service is turned online it should be set to " + PONG)
}

func TestStorageOnline(t *testing.T)  {
	ps := NewPingSummary()
	ps.TurnStorageOnline()
	assert.Equal(t, PONG, ps.Storage, "As the Storage is turned online it should be set to " + PONG)
}

func TestJsonEncoding(t *testing.T) {
	cases := dataProviderJsonEncoding()
	for _, c := range cases {
		assert.Equal(t, c.expected, jsonEncode(c.actual), "The JSON encoded struct is not matching")
	}
}

func dataProviderJsonEncoding() []testDataJson {
	psOn := NewPingSummary()
	psOn.TurnServiceOnline()
	psOn.TurnStorageOnline()
	return []testDataJson{
		{"{\"Service\":\"\",\"Storage\":\"\"}", NewPingSummary()},
		{"{\"Service\":\"PONG\",\"Storage\":\"PONG\"}", psOn},
	}
}

type testDataJson struct {
	expected string
	actual *PingSummary
}

func jsonEncode(ps *PingSummary) string {
	res, _ := json.Marshal(ps)
	return string(res)
}