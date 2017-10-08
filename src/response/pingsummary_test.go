package response

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"encoding/json"
)

func TestPingSummaryDefaultValues(t *testing.T) {
	ps := NewPingSummary()
	assert.Equal(t, "", ps.Service, "Service should return an empty string as default")
	assert.Equal(t, "", ps.Storage, "Storage should return an empty string as default")
}

func TestPingSummaryServiceOnline(t *testing.T)  {
	ps := NewPingSummary()
	ps.TurnServiceOnline()
	assert.Equal(t, PONG, ps.Service, "As the service is turned online it should be set to " + PONG)
}

func TestPingSummaryStorageOnline(t *testing.T)  {
	ps := NewPingSummary()
	ps.TurnStorageOnline()
	assert.Equal(t, PONG, ps.Storage, "As the Storage is turned online it should be set to " + PONG)
}

func TestPingSummaryJsonEncoding(t *testing.T) {
	cases := pingSummaryDataProviderJsonEncoding()
	for _, c := range cases {
		assert.Equal(t, c.expected, pingSummaryJsonEncode(c.actual), "The JSON encoded struct is not matching")
	}
}

func pingSummaryDataProviderJsonEncoding() []pingSummaryTestDataJson {
	psOn := NewPingSummary()
	psOn.TurnServiceOnline()
	psOn.TurnStorageOnline()
	return []pingSummaryTestDataJson{
		{"{\"service\":\"\",\"storage\":\"\"}", NewPingSummary()},
		{"{\"service\":\"PONG\",\"storage\":\"PONG\"}", psOn},
	}
}

type pingSummaryTestDataJson struct {
	expected string
	actual *PingSummary
}

func pingSummaryJsonEncode(ps *PingSummary) string {
	res, _ := json.Marshal(ps)
	return string(res)
}
