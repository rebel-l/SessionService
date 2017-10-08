package configuration

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConfigDefault(t *testing.T) {
	c := newConfig("")
	assert.Equal(t, newService(), c.Service, "Service should have a Service struct.")
	assert.Nil(t, c.Redis, "Redis should be nil")
}

func TestConfigJson(t *testing.T) {
	c := newConfig("")
	assert.Equal(t, "{\"Service\":{\"Port\":4000,\"LogLevel\":2},\"Redis\":null}", configJsonEncode(c), "The JSON encoded struct is not matching")
}

// TODO Test the loadFromFile method with (happy/unhappy) mocked JSON decoder, mocked File

func configJsonEncode(c *Config) string {
	res, _ := json.Marshal(c)
	return string(res)
}
