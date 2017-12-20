package session

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/pborman/uuid"
	"github.com/rebel-l/sessionservice/src/endpoint"
	"github.com/rebel-l/sessionservice/src/request"
	"github.com/rebel-l/sessionservice/src/response"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

// Put handles all PUT requests to session endpoint
type Put struct {
	session *Session
}

func NewPut(session *Session) *Put {
	put := new(Put)
	put.session = session
	return put
}

func (put *Put) Handler(res http.ResponseWriter, req *http.Request) {
	log.Info("Executing session PUT")
	msg := endpoint.NewMessage(res)

	// read request body
	requestBody, err := put.getRequestBody(req)
	if err != nil {
		log.Errorf("Parsing request body failed: %s", err)
		msg.Plain(BadRequestText, http.StatusBadRequest)
		return
	}

	// store session
	session := response.NewSession(requestBody.Id, put.session.Config.SessionLifetime)
	lifetime := time.Duration(session.Lifetime) * time.Second
	status := http.StatusOK
	var dataJson []byte
	if requestBody.Id == "" {
		log.Debugf("Create new session: %s", session.Id)
		dataJson, err = json.Marshal(requestBody.Data)
		if err != nil {
			log.Errorf("Saving Id %s failed: %s", session.Id, err)
			msg.InternalServerError(InternalServerErrorText)
			return
		}
		status = http.StatusCreated
	} else {
		log.Debugf("Update session: %s", requestBody.Id)

		// 1. load stored session
		result := put.session.Redis.Get(requestBody.Id)

		// 2. if key not found ==> return error (404)
		storageData, err := result.Result()
		if err != nil {
			log.Errorf("Session Id %s not found or has expired: %s", session.Id, err)
			msg.Plain("Session was not found or has expired.",	http.StatusNotFound)
			return
		}

		// 3. merge data with current stored
		log.Debugf("Loaded session data for %s: %s", session.Id, storageData)
		oldData := make(map[string]string)
		err = json.Unmarshal([]byte(storageData), &oldData)
		if err != nil {
			log.Errorf("Data loaded for %s can't be turned into map: %s", session.Id, err)
			msg.InternalServerError(InternalServerErrorText)
			return
		}

		requestBody.Data = put.mergeData(oldData, requestBody.Data)
		dataJson, err = json.Marshal(requestBody.Data)
		if err != nil {
			log.Errorf("Saving Id %s failed: %s", session.Id, err)
			msg.InternalServerError(InternalServerErrorText)
			return
		}
	}

	result := put.session.Redis.Set(session.Id, dataJson, lifetime)
	if result.Err() != nil {
		log.Errorf("Saving Id %s failed: %s", session.Id, result.Err().Error())
		msg.InternalServerError(InternalServerErrorText)
		return
	}
	session.Data = requestBody.Data

	for key, value := range requestBody.Data {
		log.Debugf("%s: %s", key, value)
	}

	// write request
	msg.SendJson(session, status)

	log.Info("Executing session PUT done!")
}

func (put *Put) getRequestBody(req *http.Request) (body request.Update, err error) {
	// read request body
	decoder := json.NewDecoder(req.Body)
	defer req.Body.Close()
	err = decoder.Decode(&body)
	if err != nil {
		err = errors.New(fmt.Sprintf("Unable to read request body: %s", err))
		return
	}

	err = put.validateRequestBody(&body)
	return
}

// TODO: middleware?
func (put *Put) validateRequestBody(body *request.Update) error {
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

func (put *Put) mergeData(old map[string]string, new map[string]string) map[string]string {
	result := make(map[string]string)
	for key, value := range old {
		result[key] = value
	}

	for key, value := range new {
		result[key] = value
	}

	return result
}
