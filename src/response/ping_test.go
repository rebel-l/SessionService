package response

import (
	"testing"
	"github.com/stretchr/testify/assert"
	//"encoding/json"
	"encoding/json"
)

func TestPingDefaultValues(t *testing.T) {
	p := NewPing()
	assert.Equal(t, FAILURE, p.Success, "Success should return an " + FAILURE + " string as default")
	assert.Equal(t, new(PingSummary), p.Summary, "Summary should return an empty PingSummary as default")
}

func TestPingNotify(t *testing.T) {
	cases := pingDataProviderNotify()
	for _, c := range cases {
		c.actual.Notify()
		msg := "Notify should set Success to " + c.expected
		msg += " if PingSummary.Service is " + c.actual.Summary.Service
		msg += " and PingSummary.Storage is " + c.actual.Summary.Storage
		assert.Equal(t, c.expected, c.actual.Success,  msg)
	}
}

func pingDataProviderNotify() []pingData {
	// only Service has PONG
	psServiceOnly := NewPingSummary()
	psServiceOnly.TurnServiceOnline()

	// only Storage has PONG
	psStorageOnly := NewPingSummary()
	psStorageOnly.TurnStorageOnline()

	// both Service and Storage has PONG
	psBoth := NewPingSummary()
	psBoth.TurnServiceOnline()
	psBoth.TurnStorageOnline()

	return []pingData{
		{FAILURE, NewPing()},
		{FAILURE, newPing(psServiceOnly)},
		{FAILURE, newPing(psStorageOnly)},
		{SUCCESS, newPing(psBoth)},
	}
}

func TestPingJsonEncoding(t *testing.T)  {
	cases := pingDataProviderJson()
	for _, c := range cases {
		assert.Equal(t, c.expected, pingJsonEncode(c.actual), "The JSON encoded struct is not matching")
	}
}

func pingDataProviderJson() []pingData {
	ps := NewPingSummary()
	ps.TurnServiceOnline()
	ps.TurnStorageOnline()

	p := newPing(ps)
	p.Notify()

	return []pingData {
		{"{\"success\":\"FAIL\",\"summary\":{\"service\":\"\",\"storage\":\"\"}}", NewPing()},
		{"{\"success\":\"OK\",\"summary\":{\"service\":\"PONG\",\"storage\":\"PONG\"}}", p},
	}
}

func pingJsonEncode(ps *Ping) string {
	res, _ := json.Marshal(ps)
	return string(res)
}

type pingData struct {
	expected string
	actual *Ping
}
