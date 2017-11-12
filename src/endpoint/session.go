package endpoint

import (
	"encoding/json"
	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
	"github.com/rebel-l/sessionservice/src/authentication"
	"github.com/rebel-l/sessionservice/src/request"
	"github.com/rebel-l/sessionservice/src/response"
	log "github.com/sirupsen/logrus"
	"net/http"
)

// Session handles the session endpoints
type Session struct {
	redis *redis.Client
	middleware *authentication.Authentification
}

// InitSession initializes the session endpoints
func InitSession(redisClient *redis.Client, router *mux.Router, middleware *authentication.Authentification) {
	log.Debug("Session endpoint: Init ...")

	s := new(Session)
	s.redis = redisClient
	s.middleware = middleware
	router.Handle("/session/", s.handlerFactory(http.MethodPut)).Methods(http.MethodPut)

	log.Debug("Ping endpoint: initialized!")
}

func (s *Session) handlePut(res http.ResponseWriter, req *http.Request) {
	log.Println("Executing session PUT")

	// read request body
	decoder := json.NewDecoder(req.Body)
	defer req.Body.Close()
	var data request.Update
	err := decoder.Decode(&data)
	if err != nil {
		log.Errorf("Unable to read request body: %s", err)
	}

	log.Debugf("Id to update: %s", data.Id)
	for key, value := range data.Data {
		log.Debugf("%s: %s", key, value)
	}

	// write request
	res.Header().Set(contentHeader, contentType)
	res.WriteHeader(http.StatusOK)
	session := response.NewSession("",0)
	err = json.NewEncoder(res).Encode(session)
	if err != nil {
		log.Errorf("Wasn't able to write body: %s", err)
	}

	log.Info("Executing session PUT done!")
}

func (s *Session) handlerFactory(method string) http.Handler {
	var handler func(http.ResponseWriter, *http.Request)
	switch method {
		case http.MethodPut:
			handler = s.handlePut
		default:
			log.Panicf("Method %s is not implemented", method)
			return nil
	}

	return s.middleware.Middleware(http.HandlerFunc(handler))
}
