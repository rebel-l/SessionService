package storage

import (
	"time"
)

type Handler interface {
	Get(id string) (string, error)
	Ping() (string, error)
	Set(key string, value interface{}, lifetime time.Duration) error
}
