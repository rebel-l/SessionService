package session

import (
	"encoding/json"
	"errors"
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

// NewSession creates a new session object
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
	router.Handle("/session/", s.handlerFactory(http.MethodGet)).Methods(http.MethodGet)

	log.Debug("Session endpoint: initialized!")
}

// handlerFactory is a factory to create the endpoint handlers
func (s *Session) handlerFactory(method string) http.Handler {
	var handler func(http.ResponseWriter, *http.Request)
	switch method {
		case http.MethodPut:
			put := NewPut(s)
			handler = put.Handler
		case http.MethodGet:
			get := NewGet(s)
			handler = get.Handler
		default:
			log.Panicf("Method %s is not implemented", method)
			return nil
	}

	return s.Authentication.Middleware(http.HandlerFunc(handler))
}

func (s *Session) loadData(id string) (data map[string]string, err error, code int) {
	// 1. load stored session
	storageData, err := s.Storage.Get(id)

	// 2. if key not found ==> respond error (404)
	if err != nil {
		log.Errorf("Session Id %s not found or has expired: %s", id, err)
		code = http.StatusNotFound
		err = errors.New(SessionNotFoundText)
		return
	}

	log.Debugf("Loaded session data for %s: %s", id, storageData)
	data = make(map[string]string)
	err = json.Unmarshal([]byte(storageData), &data)
	if err != nil {
		log.Errorf("Data loaded for %s can't be turned into map: %s", id, err)
		code = http.StatusInternalServerError
		err = errors.New(InternalServerErrorText)
		return
	}

	code = http.StatusOK
	return
}
