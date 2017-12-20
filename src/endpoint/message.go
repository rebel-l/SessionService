package endpoint

import (
	log "github.com/sirupsen/logrus"
	"net/http"
	"encoding/json"
)

type Message struct {
	res http.ResponseWriter
}

func NewMessage(res http.ResponseWriter) Message {
	return Message{res}
}

func (m *Message) Plain(msg string, code int) {
	m.res.Header().Set(ContentHeader, ContentTypePlain)
	m.res.WriteHeader(code)
	i,_ := m.res.Write([]byte(msg))
	if i < 1 {
		log.Errorf("Wasn't able to write body: %d", i)
	}
}

func (m *Message) InternalServerError(msg string) {
	m.Plain(
		msg,
		http.StatusInternalServerError,
	)
}

func (m *Message) SendJson(value interface{}, code int) {
	m.res.Header().Set(ContentHeader, ContentTypeJson)
	m.res.WriteHeader(code)
	err := json.NewEncoder(m.res).Encode(value)
	if err != nil {
		log.Errorf("Wasn't able to write body: %s", err)
	}
}