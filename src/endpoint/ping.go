package endpoint

import (
	"encoding/json"
	"github.com/rebel-l/sessionservice/src/response"
	log "github.com/sirupsen/logrus"
	"net/http"
)

const contentHeader = "Content-Type"
const contentType = "application/json"

type Ping struct {
	response *response.Ping
	observer []response.Observer
}

func InitPing() {
	log.Debug("Ping endpoint: Init ...")

	p := new(Ping)
	p.response = response.NewPing()
	p.observer = append(p.observer, p.response)
	http.HandleFunc("/ping/", p.handler)

	log.Debug("Ping endpoint: initialized!")
}

func (p *Ping) handler(res http.ResponseWriter, req *http.Request) {
	log.Debug("Ping: request received ...")

	// do the checks
	p.checkService()
	p.checkStorage()

	// send response
	p.send(res)

	log.Debug("Ping: response send!")
}

func (p *Ping) notify()  {
	for _, o := range p.observer {
		o.Notify()
	}
}

func (p *Ping) checkService() {
	p.response.Summary.TurnServiceOnline()
	p.notify()
}

func (p *Ping) checkStorage() {
	// TODO: implement redis ping here ...
	p.response.Summary.TurnStorageOnline()
	p.notify()
}

func (p *Ping) send(res http.ResponseWriter)  {
	res.Header().Set(contentHeader, contentType)
	if p.response.Success != response.SUCCESS {
		res.WriteHeader(http.StatusServiceUnavailable)
	}
	json.NewEncoder(res).Encode(p.response)
}
