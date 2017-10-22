package configuration

import (
	"github.com/stretchr/testify/assert"
	"testing"
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

// TODO: test parameters set the config correct

// TODO: test that parameters overwrite config file
