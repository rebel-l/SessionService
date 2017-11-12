package endpoint

import (
	"encoding/json"
	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
	"github.com/rebel-l/sessionservice/src/response"
	log "github.com/sirupsen/logrus"
	"net/http"
	"sync"
	"time"
)

// Ping handles the ping endpoint
type Ping struct {
	response *response.Ping
	observer []response.Observer
	redisClient *redis.Client
}

func InitPing(redisClient *redis.Client, r *mux.Router) {
	log.Debug("Ping endpoint: Init ...")

	p := new(Ping)
	p.response = response.NewPing()
	p.observer = append(p.observer, p.response)
	p.redisClient = redisClient
	r.HandleFunc("/ping/", p.handler).Methods(http.MethodGet)

	log.Debug("Ping endpoint: initialized!")
}

func (p *Ping) handler(res http.ResponseWriter, req *http.Request) {
	log.Debug("Ping: request received ...")
	start := time.Now()
	wg := new(sync.WaitGroup)
	wg.Add(2)

	// do the checks
	go p.checkService(wg)
	go p.checkStorage(wg)

	// send response
	wg.Wait()
	p.send(res)
	stop := time.Now()
	log.Infof("Ping: response send! Duration: %s", (stop.Sub(start)).String())
}

func (p *Ping) notify()  {
	for _, o := range p.observer {
		o.Notify()
	}
}

func (p *Ping) checkService(wg *sync.WaitGroup) {
	p.response.Summary.TurnServiceOnline()
	p.notify()
	wg.Done()
}

func (p *Ping) checkStorage(wg *sync.WaitGroup) {
	pong, err := p.redisClient.Ping().Result()
	if err != nil {
		log.Errorf("Redis storage is not available: %s", err)
	} else {
		log.Debugf("Redis Ping responded with %s", pong)
		p.response.Summary.TurnStorageOnline()
	}
	p.notify()
	wg.Done()
}

func (p *Ping) send(res http.ResponseWriter)  {
	res.Header().Set(contentHeader, contentType)
	if p.response.Success != response.SUCCESS {
		res.WriteHeader(http.StatusServiceUnavailable)
	}
	json.NewEncoder(res).Encode(p.response)
}
