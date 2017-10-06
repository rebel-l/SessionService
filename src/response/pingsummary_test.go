package response

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestDefaultValues(t *testing.T) {
	ps := NewPingSummary()
	assert.Equal(t, "", ps.Service(), "Service should return an empty string as default")
	assert.Equal(t, "", ps.Storage(), "Storage should return an empty string as default")
}

func TestServiceOnline(t *testing.T)  {
	ps := NewPingSummary()
	ps.ServiceOnline()
	assert.Equal(t, PONG, ps.Service(), "As the service is turned online it should be set to " + PONG)
}

func TestStorageOnline(t *testing.T)  {
	ps := NewPingSummary()
	ps.StorageOnline()
	assert.Equal(t, PONG, ps.Storage(), "As the storage is turned online it should be set to " + PONG)
}
