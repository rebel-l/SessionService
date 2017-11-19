package configuration

import (
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
	"os"
)

type ParserMock struct {
	counter int
}

func (pm *ParserMock) printVersion() {
	pm.counter++
}

func NewParserMock() *ParserMock {
	pm := new(ParserMock)
	pm.counter = 0
	return pm
}

func TestParserFlagVersionUnhappy(t *testing.T) {
	pm := NewParserMock()
	p := GetParser()
	p.printVersion = pm.printVersion
	p.Parse()
	assert.Equal(t, 0, pm.counter, "If version flag is NOT set, printVersion() method should have been NOT called")
}

func TestParserFlagVersionHappy(t *testing.T) {
	// Setup
	pm := NewParserMock()
	p := GetParser()
	p.printVersion = pm.printVersion
	p.version = true

	// Test
	p.Parse()
	assert.Equal(t, 1, pm.counter, "If version flag is set, printVersion() method should have been called")
}

func TestParserParametersSetConfig(t *testing.T) {
	// Setup
	p := GetParser()
	p.servicePort = 1000
	p.serviceLogLevel = 6
	p.serviceSessionLifetime = 100
	p.redisAddr = "redisUrl"
	p.redisDb = 21
	p.redisPassword = "secret"

	// Test
	config := p.Parse()
	assert.Equal(t, 1000, config.Service.Port, "Service.Port couldn't be set by parser (parameter)")
	assert.Equal(t, log.Level(6), config.Service.LogLevel, "Service.LogLevel couldn't be set by parser (parameter)")
	assert.Equal(t, 100, config.Service.SessionLifetime, "Service.SessionLifetime couldn't be set by parser (parameter)")
	assert.Equal(t, "redisUrl", config.Redis.Addr, "Redis.Addr couldn't be set by parser (parameter)")
	assert.Equal(t, 21, config.Redis.DB, "Redis.DB couldn't be set by parser (parameter)")
	assert.Equal(t, "secret", config.Redis.Password, "Redis.Password couldn't be set by parser (parameter)")
}

func TestParserParametersOverwriteConfigFromFile(t *testing.T) {
	/**
	 * Part 1: ensure file is read
	 */
	// Setup
	p := GetParser()
	p.servicePort = 0
	p.serviceLogLevel = 0
	p.serviceSessionLifetime = 0
	p.redisAddr = ""
	p.redisDb = 0
	p.redisPassword = ""
	p.filename = os.Getenv("GOPATH")
	p.filename += "/src/github.com/rebel-l/sessionservice/testFixtures/configuration/parser/config.json"

	// Test
	config := p.Parse()
	assert.Equal(t, 333, config.Service.Port, "Service.Port was not loaded from file")
	assert.Equal(t, log.Level(5), config.Service.LogLevel, "Service.LogLevel was not loaded from file")
	assert.Equal(t, 3600, config.Service.SessionLifetime, "Service.SessionLifetime was not loaded from file")
	assert.Equal(t, "redis:6379", config.Redis.Addr, "Redis.Addr was not loaded from file")
	assert.Equal(t, 11, config.Redis.DB, "Redis.DB was not loaded from file")
	assert.Equal(t, "1234", config.Redis.Password, "Redis.Password was not loaded from file")

	/**
	 * Part 2: ensure cli parameters overwrite config from file
	 */
	 // Setup
	p.servicePort = 1000
	p.serviceLogLevel = 6
	p.serviceSessionLifetime = 3000
	p.redisAddr = "redisUrl"
	p.redisDb = 21
	p.redisPassword = "secret"

	// Test
	config = p.Parse()
	assert.Equal(t, 1000, config.Service.Port, "Service.Port couldn't be set by parser (parameter)")
	assert.Equal(t, log.Level(6), config.Service.LogLevel, "Service.LogLevel couldn't be set by parser (parameter)")
	assert.Equal(t, 3000, config.Service.SessionLifetime, "Service.SessionLifetime couldn't be set by parser (parameter)")
	assert.Equal(t, "redisUrl", config.Redis.Addr, "Redis.Addr couldn't be set by parser (parameter)")
	assert.Equal(t, 21, config.Redis.DB, "Redis.DB couldn't be set by parser (parameter)")
	assert.Equal(t, "secret", config.Redis.Password, "Redis.Password couldn't be set by parser (parameter)")
}
