package storage

import (
	"github.com/go-redis/redis"
	"time"
)

type Redis struct {
	client *redis.Client
}

func (r *Redis) Get(id string) (string, error) {
	return r.client.Get(id).Result()
}

func NewRedis(client *redis.Client) *Redis {
	r := new(Redis)
	r.client = client
	return r
}

func (r *Redis) Ping() (string, error) {
	return r.client.Ping().Result()
}

func (r *Redis) Set(key string, value interface{}, lifetime time.Duration) error {
	return r.client.Set(key, value, lifetime).Err()
}
