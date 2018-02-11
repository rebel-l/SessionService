package session

import (
	"errors"
	"github.com/pborman/uuid"
	"github.com/rebel-l/sessionservice/src/endpoint"
	"github.com/rebel-l/sessionservice/src/request"
	"github.com/rebel-l/sessionservice/src/response"
	log "github.com/sirupsen/logrus"
	"net/http"
	"fmt"
	"strconv"
)

const ValidateErrorMessage = "request parameters validation failed ==> no UUID provided"

type Get struct {
	session *Session
}

func NewGet(session *Session) *Get {
	get := new(Get)
	get.session = session
	return get
}

func (get *Get) Handler(res http.ResponseWriter, req *http.Request) {
	log.Info("Executing session GET")
	msg := endpoint.NewMessage(res)

	// read params
	params := get.parseParams(req)

	// validate params
	err := get.validateRequestParams(params)
	if err != nil {
		log.Errorf("Validating request parameters failed: %s", err)
		msg.Plain(BadRequestTextGet, http.StatusBadRequest)
		return
	}

	// load session
	sessionResponse := response.NewSession(params.Id, get.session.Config.SessionLifetime)
	code := http.StatusOK
	sessionResponse.Data, err, code = get.session.loadData(params.Id)
	if err != nil {
		msg.Plain(err.Error(), code)
		return
	}

	// set new lifetime
	// TODO: implement renew of lifetime

	// write response
	msg.SendJson(sessionResponse, code)

	log.Info("Executing session GET done!")
}

func (get *Get) parseParams(req *http.Request) *request.Read {
	params := new(request.Read)
	query := req.URL.Query()
	fmt.Printf("params: %#v\n", query)
	params.Id = query.Get("id")
	log.Debugf("UUID parsed: %s", params.Id)
	regenerateId, err := strconv.ParseBool(query.Get("regenerateId"))
	if err != nil {
		regenerateId = false
	}
	params.RegenerateId = regenerateId
	return params
}

func (get *Get) validateRequestParams(params *request.Read) error {
	// Id must be not empty and uuid
	if params.Id != "" {
		params.Id = uuid.Parse(params.Id).String()

		if len(params.Id) != UUIDLENGTH {
			return errors.New(ValidateErrorMessage)
		}
	} else {
		return errors.New(ValidateErrorMessage)
	}

	return nil
}
