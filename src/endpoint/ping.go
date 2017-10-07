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
}

func InitPing() {
	log.Debug("Ping endpoint: Init ...")

	p := new(Ping)
	p.response = response.NewPing()
	http.HandleFunc("/ping/", p.handler)

	log.Debug("Ping endpoint: initialized!")
}

func (p *Ping) handler(res http.ResponseWriter, req *http.Request) {
	log.Debug("Ping: request received ...")

	res.Header().Set(contentHeader, contentType)
	json.NewEncoder(res).Encode(p.response)

	log.Debug("Ping: response send!")
}