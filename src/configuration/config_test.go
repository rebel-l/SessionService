package configuration

import (
	//"encoding/json"
	"github.com/go-redis/redis"
	"github.com/stretchr/testify/assert"
	"testing"
	"os"
	"errors"
)

func TestConfigDefault(t *testing.T) {
	c := newConfig("")
	assert.Equal(t, newService(), c.Service, "Service should have a Service struct.")
	assert.Equal(t, new(redis.Options), c.Redis, "Redis should be nil")
}

type fakeOs struct {
	file *os.File
	err error
	openCounter int
}

func (f *fakeOs) Open(name string) (file *os.File, err error) {
	f.openCounter++
	file = f.file
	err = f.err
	return
}

func TestConfigLoadFileUnhappy(t *testing.T) {
	f := &fakeOs{ file: new(os.File), err: errors.New("File doesn't exist"), openCounter: 0}
	c := newConfig("")
	c.openFile = f.Open
	err := c.loadFromFile("notexisting")
	assert.Equal(t, 1, f.openCounter, "Open() method should have called once")
	assert.Equal(t, f.err, err, "Expected that method is returning the error message")
}

// TODO Test the loadFromFile method with (happy/unhappy) mocked JSON decoder, mocked File


