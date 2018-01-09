package storage

import (
	"github.com/go-redis/redis"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStorageNewRedis(t *testing.T) {
	r := redis.NewClient(&redis.Options{})
	s := NewRedis(r)
	assert.Equal(t, r, s.client, "Redis client was not set")
}
