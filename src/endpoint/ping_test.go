package endpoint

import (
	"github.com/go-redis/redis"
	"github.com/rebel-l/sessionservice/src/response"
	"github.com/stretchr/testify/assert"
	"os"
	"sync"
	"testing"
)

var redisAddr = os.Getenv("REDISADDR")

type ObserverMock struct {
	counter int
}

func (o *ObserverMock) Notify() {
	o.counter++
}

func getPingWithObserverMock() (p *Ping, o *ObserverMock) {
	p = new(Ping)
	o = new(ObserverMock)
	p.observer = append(p.observer, o)
	return
}

func getPingWithObserverMockAndResponse() (p *Ping, o *ObserverMock) {
	p, o = getPingWithObserverMock()
	p.response = response.NewPing()
	p.observer = append(p.observer, p.response)
	return
}

func getPingWithObserverMockComplete() (p *Ping, o *ObserverMock) {
	p, o = getPingWithObserverMockAndResponse()
	redisOptions := new(redis.Options)
	if redisAddr == "" {
		redisAddr = "redis:6379"
	}
	redisOptions.Addr = redisAddr
	p.redisClient = redis.NewClient(redisOptions)
	return
}

func TestPingNotifyHappy(t *testing.T) {
	p, o := getPingWithObserverMock()
	p.notify()
	assert.Equal(t, 1, o.counter, "Notify() method of observers should have been called")
}

func TestPingCheckServicePartlyHappy(t *testing.T) {
	p, o := getPingWithObserverMockAndResponse()
	wg := new(sync.WaitGroup)
	wg.Add(1)
	p.checkService(wg)
	assert.Equal(t, 1, o.counter, "Observers should be notified")
	assert.Equal(t, response.FAILURE, p.response.Success)
	assert.Equal(t, response.PONG, p.response.Summary.Service)
}

func TestPingCheckStoragePartlyHappy(t *testing.T) {
	p, o := getPingWithObserverMockComplete()
	wg := new(sync.WaitGroup)
	wg.Add(1)
	p.checkStorage(wg)
	assert.Equal(t, 1, o.counter, "Observers should be notified")
	assert.Equal(t, response.FAILURE, p.response.Success)
	assert.Equal(t, response.PONG, p.response.Summary.Storage)
}

func TestPingCheckFullHappy(t *testing.T) {
	p, o := getPingWithObserverMockComplete()
	wg := new(sync.WaitGroup)
	wg.Add(2)
	p.checkService(wg)
	p.checkStorage(wg)
	assert.Equal(t, 2, o.counter, "Observers should be notified for each check")
	assert.Equal(t, response.SUCCESS, p.response.Success)
	assert.Equal(t, response.PONG, p.response.Summary.Service)
	assert.Equal(t, response.PONG, p.response.Summary.Storage)
}
