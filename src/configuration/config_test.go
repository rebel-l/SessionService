package configuration

import (
	"errors"
	"github.com/go-redis/redis"
	"github.com/stretchr/testify/assert"
	log "github.com/sirupsen/logrus"
	"os"
	"testing"
)

func TestInit(t *testing.T) {
	log.SetLevel(0)
}

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

func TestConfigLoadFileHappy(t *testing.T) {
	filename := os.Getenv("GOPATH")
	filename += "/src/github.com/rebel-l/sessionservice/testFixtures/configuration/config/happy.json"
	c := newConfig(filename)
	assert.Equal(t, 333, c.Service.Port, "Port was not correctly loaded from file")
	assert.Equal(t, 5, int(c.Service.LogLevel), "Loglevel was not correctly loaded from file")
	assert.Equal(t, "redis:6379", c.Redis.Addr, "Redis host was not correctly loaded from file")
	assert.Equal(t, "1234", c.Redis.Password, "Redis password was not correctly loaded from file")
	assert.Equal(t, 0, c.Redis.DB, "Redis DB was not correctly loaded from file")
	assert.Equal(t, "secretapitoken", c.AccountList["MyApp"].ApiKey, "AccountList was not loaded correct")
	assert.Equal(t, "secretapitoken2", c.AccountList["OtherApp"].ApiKey, "AccountList was not loaded correct")
}

func TestConfigDecodeUnhappy(t *testing.T) {
	filename := os.Getenv("GOPATH")
	filename += "/src/github.com/rebel-l/sessionservice/testFixtures/configuration/config/unhappy.json"
	c := newConfig("")
	err := c.loadFromFile(filename)
	assert.Equal(t, "EOF", err.Error(), "Decode should fail on malformed JSON")
}
