package storage

import (
	"github.com/stretchr/testify/mock"
	"time"
)

type HandlerMock struct {
	mock.Mock
}

func (m *HandlerMock) Get(id string) (string, error) {
	args := m.Called(id)
	return args.String(0), args.Error(1)
}

func (m *HandlerMock) Ping() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

func (m *HandlerMock) Set(key string, value interface{}, lifetime time.Duration) error {
	args := m.Called(key, value, lifetime)
	return args.Error(0)
}
