package session

import (
	"encoding/json"
	"github.com/rebel-l/sessionservice/src/authentication"
	"github.com/rebel-l/sessionservice/src/configuration"
	"github.com/rebel-l/sessionservice/src/response"
	"github.com/rebel-l/sessionservice/src/storage"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
	"errors"
	"net/http"
)

func TestEndpointSessionPutNewPut(t *testing.T) {
	storage := new(storage.HandlerMock)
	session := getSessionMock(storage)
	put := NewPut(session)
	assert.Equal(t,session, put.session, "Expected session struct added to put struct")
	storage.AssertExpectations(t)
}

func TestEndpointSessionPutStoreDataHappy(t *testing.T) {
	// setup fixtures
	id := "sessionId"
	data := make(map[string]string)
	data["Json"] = "yes"
	dataJson, err := json.Marshal(data)
	if err != nil {
		t.Error("Error on converting to JSON")
	}
	res := &response.Session{Id: id, Lifetime: 123}
	lifetime := time.Duration(res.Lifetime)

	// setup mock
	storage := new(storage.HandlerMock)
	storage.On("Set", id, dataJson, lifetime * time.Second).Return(nil)

	// do the test
	session := getSessionMock(storage)
	put := NewPut(session)
	assert.Empty(t, res.Data, "Data needs to be empty before storing data")
	assert.Nil(t, put.storeData(res, data), "Wasn't able to store data")
	assert.Equal(t, data, res.Data, "Data was not added to response")
	storage.AssertExpectations(t)
}

func TestEndpointSessionPutStoreDataUnhappy(t *testing.T) {
	// setup fixtures
	id := "sessionIdNew"
	data := make(map[string]string)
	data["key"] = "value"
	dataJson, err := json.Marshal(data)
	if err != nil {
		t.Error("Error on converting to JSON")
	}
	res := &response.Session{Id: id, Lifetime: 123}
	lifetime := time.Duration(res.Lifetime)
	errMsg := "Saving Id sessionIdNew failed: Failing storing data"

	// setup mock
	storage := new(storage.HandlerMock)
	storage.On("Set", id, dataJson, lifetime * time.Second).Return(errors.New("Failing storing data"))

	// do the test
	session := getSessionMock(storage)
	put := NewPut(session)
	assert.Empty(t, res.Data, "Data needs to be empty before storing data")
	assert.Equal(t, errMsg, put.storeData(res, data).Error(), "Error should be returned after failing")
	assert.Empty(t, res.Data, "Data in response should be still empty after failing")
	storage.AssertExpectations(t)
}

func TestEndpointSessionPutLoadDataHappy(t *testing.T) {
	// setup fixtures
	id := "existingId"
	data := make(map[string]string)
	data["someKey"] = "a boring value"
	dataJson, err := json.Marshal(data)
	if err != nil {
		t.Error("Error on converting to JSON")
	}

	// setup mock
	storage := new(storage.HandlerMock)
	storage.On("Get", id).Return(string(dataJson), nil)

	// do the test
	session := getSessionMock(storage)
	put := NewPut(session)
	result, err, code := put.loadData(id)
	assert.Equal(t, data, result, "Data was not loaded correct")
	assert.Nil(t, err, "There should be not error on happy path")
	assert.Equal(t, http.StatusOK, code, "The http code should show a success")
	storage.AssertExpectations(t)
}

func TestEndpointSessionPutLoadDataUnhappy(t *testing.T) {
	// setup fixtures
	id := "existingId"
	data := make(map[string]string)
	data["someKey"] = "a boring value"
	errMsg := "Session was not found or has expired."

	// setup mock
	storage := new(storage.HandlerMock)
	storage.On("Get", id).Return("", errors.New("Failing loading data"))

	// do the test
	session := getSessionMock(storage)
	put := NewPut(session)
	result, err, code := put.loadData(id)
	assert.Nil(t, result, "Data should be not returned on fail")
	assert.Equal(t, errMsg, err.Error(), "There should be an error on fail")
	assert.Equal(t, http.StatusNotFound, code, "The http code should show a not found")
	storage.AssertExpectations(t)
}

/**
ToDo: missing tests ...
	Handler()
	getRequestBody()
	validateRequestBody()
	mergeData()
 */

func getSessionMock(storage storage.Handler) *Session {
	auth := &authentication.Authentication{}
	conf := &configuration.Service{}
	return NewSession(storage, auth, conf)
}
