package endpoint

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
	"github.com/pborman/uuid"
	"github.com/rebel-l/sessionservice/src/authentication"
	"github.com/rebel-l/sessionservice/src/configuration"
	"github.com/rebel-l/sessionservice/src/request"
	"github.com/rebel-l/sessionservice/src/response"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

const (
	putBadRequestText = "A body needs to be send as JSON and it needs at least a data field"
)

// Session handles the session endpoints
type Session struct {
	redis *redis.Client
	middleware *authentication.Authentification
	config *configuration.Service
}

// InitSession initializes the session endpoints
func InitSession(redisClient *redis.Client, router *mux.Router, middleware *authentication.Authentification, config *configuration.Service) {
	log.Debug("Session endpoint: Init ...")

	// init Session struct
	s := new(Session)
	s.redis = redisClient
	s.middleware = middleware
	s.config = config

	// register handler
	router.Handle("/session/", s.handlerFactory(http.MethodPut)).Methods(http.MethodPut)

	log.Debug("Session endpoint: initialized!")
}

func (s *Session) handlePut(res http.ResponseWriter, req *http.Request) {
	log.Info("Executing session PUT")

	// read request body
	requestBody, err := s.getRequestBody(req)
	if err != nil {
		log.Errorf("Parsing request body failed: %s", err)
		s.sendPlain(putBadRequestText, res, http.StatusBadRequest)
		return
	}

	// store session
	session := response.NewSession(requestBody.Id, s.config.SessionLifetime)
	lifetime := time.Duration(session.Lifetime) * time.Second
	status := http.StatusOK
	var dataJson []byte
	if requestBody.Id == "" {
		log.Debugf("Create new session: %s", session.Id)
		dataJson, err = json.Marshal(requestBody.Data)
		if err != nil {
			log.Errorf("Saving Id %s failed: %s", session.Id, err)
			s.sendInternalServerError(res)
			return
		}
		status = http.StatusCreated
	} else {
		log.Debugf("Update session: %s", requestBody.Id)

		// 1. load stored session
		result := s.redis.Get(requestBody.Id)

		// 2. if key not found ==> return error (404)
		storageData, err := result.Result()
		if err != nil {
			log.Errorf("Session Id %s not found or has expired: %s", session.Id, err)
			s.sendPlain(
				"Session was not found or has expired.",
				res,
				http.StatusNotFound,
			)
			return
		}

		// 3. merge data with current stored
		log.Debugf("Loaded session data for %s: %s", session.Id, storageData)
		oldData := make(map[string]string)
		err = json.Unmarshal([]byte(storageData), &oldData)
		if err != nil {
			log.Errorf("Data loaded for %s can't be turned into map: %s", session.Id, err)
			s.sendInternalServerError(res)
			return
		}

		requestBody.Data = s.mergeData(oldData, requestBody.Data)
		dataJson, err = json.Marshal(requestBody.Data)
		if err != nil {
			log.Errorf("Saving Id %s failed: %s", session.Id, err)
			s.sendInternalServerError(res)
			return
		}
	}

	result := s.redis.Set(session.Id, dataJson, lifetime)
	if result.Err() != nil {
		log.Errorf("Saving Id %s failed: %s", session.Id, result.Err().Error())
		s.sendInternalServerError(res)
		return
	}
	session.Data = requestBody.Data

	for key, value := range requestBody.Data {
		log.Debugf("%s: %s", key, value)
	}

	// write request
	res.Header().Set(contentHeader, contentTypeJson)
	res.WriteHeader(status)
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

func (s *Session) getRequestBody(req *http.Request) (body request.Update, err error) {
	// read request body
	decoder := json.NewDecoder(req.Body)
	defer req.Body.Close()
	err = decoder.Decode(&body)
	if err != nil {
		err = errors.New(fmt.Sprintf("Unable to read request body: %s", err))
		return
	}

	err = s.validateRequestBody(&body)
	return
}

func (s *Session) validateRequestBody(body *request.Update) error {
	// Id must be uuid or empty string
	if body.Id != "" {
		body.Id = uuid.Parse(body.Id).String()
		log.Debugf("UUID parsed: %s", body.Id)
		if len(body.Id) != UUIDLENGTH {
			return errors.New("request body validation failed ==> wrong UUID provided")
		}
	}

	// data field must have entries
	if len(body.Data) < 1 {
		return errors.New("request body validation failed ==> no data field received")
	}

	return nil
}

func (s *Session) sendPlain(msg string, res http.ResponseWriter, code int) {
	res.Header().Set(contentHeader, contentTypePlain)
	res.WriteHeader(code)
	i,_ := res.Write([]byte(msg))
	if i < 1 {
		log.Errorf("Wasn't able to write body: %d", i)
	}
}

func (s *Session) sendInternalServerError(res http.ResponseWriter) {
	s.sendPlain(
		"Session could not be stored because of internal error. Contact administrator or retry it later.",
		res,
		http.StatusInternalServerError,
	)
}

func (s *Session) mergeData(old map[string]string, new map[string]string) map[string]string {
	result := make(map[string]string)
	for key, value := range old {
		result[key] = value
	}

	for key, value := range new {
		result[key] = value
	}

	return result
}
