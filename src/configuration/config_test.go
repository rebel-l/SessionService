package configuration

import (
	//"encoding/json"
	"github.com/go-redis/redis"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConfigDefault(t *testing.T) {
	c := newConfig("")
	assert.Equal(t, newService(), c.Service, "Service should have a Service struct.")
	assert.Equal(t, new(redis.Options), c.Redis, "Redis should be nil")
}


// TODO Test the loadFromFile method with (happy/unhappy) mocked JSON decoder, mocked File


