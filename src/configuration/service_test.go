package configuration

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

func TestServiceDefault(t *testing.T)  {
	s := newService()
	assert.Equal(t, ServiceDefaultPort, s.Port, "The default Port of the configuration should be " + strconv.Itoa(ServiceDefaultPort))
	assert.Equal(t, ServiceDefaultLogLevel, s.LogLevel, "The default LogLevel of the configuration should be " + ServiceDefaultLogLevel.String())
}

func TestServiceChanges(t *testing.T)  {
	s := newService()
	s.Port = 666
	s.LogLevel = 5
	assert.Equal(t, 666, s.Port, "The Port can't be changed")
	assert.Equal(t, 5, int(s.LogLevel), "The LogLevel can't be changed")
}

func TestServiceJson(t *testing.T)  {
	s := newService()
	assert.Equal(t, "{\"Port\":4000,\"LogLevel\":2}", serviceJsonEncode(s), "The JSON encoded struct is not matching")
}

func serviceJsonEncode(s *Service) string {
	res, _ := json.Marshal(s)
	return string(res)
}
