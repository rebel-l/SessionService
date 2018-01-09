package session

import (
	"github.com/gorilla/mux"
	"github.com/rebel-l/sessionservice/src/authentication"
	"github.com/rebel-l/sessionservice/src/configuration"
	"github.com/rebel-l/sessionservice/src/storage"
	log "github.com/sirupsen/logrus"
	"net/http"
)

// Session handles the session endpoints
type Session struct {
	Storage        storage.Handler
	Authentication *authentication.Authentication
	Config         *configuration.Service
}

func NewSession(
	storage storage.Handler,
	authentication *authentication.Authentication,
	config *configuration.Service) *Session {
	s := new(Session)
	s.Storage = storage
	s.Authentication = authentication
	s.Config = config
	return s
}

// Init initializes the session endpoints
func (s *Session) Init(router *mux.Router) {
	log.Debug("Session endpoint: Init ...")

	router.Handle("/session/", s.handlerFactory(http.MethodPut)).Methods(http.MethodPut)

	log.Debug("Session endpoint: initialized!")
}

func (s *Session) handlerFactory(method string) http.Handler {
	var handler func(http.ResponseWriter, *http.Request)
	switch method {
		case http.MethodPut:
			put := NewPut(s)
			handler = put.Handler
		default:
			log.Panicf("Method %s is not implemented", method)
			return nil
	}

	return s.Authentication.Middleware(http.HandlerFunc(handler))
}
