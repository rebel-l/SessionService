package session

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/rebel-l/sessionservice/src/authentication"
	"github.com/rebel-l/sessionservice/src/configuration"
	"github.com/rebel-l/sessionservice/src/request"
	"github.com/rebel-l/sessionservice/src/response"
	"github.com/rebel-l/sessionservice/src/storage"
	"github.com/rebel-l/sessionservice/src/utils/testify"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
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

func TestEndpointSessionPutGetRequestBodyHappy(t *testing.T) {
	// setup
	requestBodyRaw, err := http.NewRequest("POST", "/session", strings.NewReader(`{"id": "myId","data":{"key":"value"}}`))
	if err != nil {
		t.Fatal(err)
	}

	storage := new(storage.HandlerMock)
	session := getSessionMock(storage)
	put := NewPut(session)

	// do the test
	body, err := put.getRequestBody(requestBodyRaw)
	assert.Nil(t, err, "getRequestBody should not cause an error")
	assert.Equal(t, "myId", body.Id, "Id was not read correct")
	assert.Equal(t, "value", body.Data["key"], "Data was parsed wrong")
}

func TestEndpointSessionPutGetRequestUnhappy(t *testing.T) {
	// setup
	requestBodyRaw, err := http.NewRequest("POST", "/session", strings.NewReader("not a JSON"))
	if err != nil {
		t.Fatal(err)
	}

	storage := new(storage.HandlerMock)
	session := getSessionMock(storage)
	put := NewPut(session)

	// do the test
	body, err := put.getRequestBody(requestBodyRaw)
	assert.True(t, strings.Contains(err.Error(), "Unable to read request body:"),"getRequestBody should cause an error")
	assert.Empty(t, body.Id)
	assert.Empty(t, body.Data)
}

func TestEndpointSessionPutValidateRequestBodyHappy(t *testing.T) {
	// setup
	body := make([]request.Update, 2)
	body[0].Data = make(map[string]string)
	body[0].Data["key"] = "value"

	body[1].Id = "8d9af075-1aa6-46c0-913d-ff42f22ca307"
	body[1].Data = body[0].Data

	storage := new(storage.HandlerMock)
	session := getSessionMock(storage)
	put := NewPut(session)

	// do the test
	for _, v := range body {
		assert.Nil(t, put.validateRequestBody(&v))
	}
}

func TestEndpointSessionPutValidateRequestBodyUnappy(t *testing.T) {
	// setup
	cases := make([]dataProviderValidateUnhappy, 2)

	// case 1: empty id, empty data ==> first object is empty
	cases[0].body = new(request.Update)
	cases[0].err = "request body validation failed ==> no data field received"

	// case 2: wrong formated UUID
	cases[1].body = new(request.Update)
	cases[1].body.Id = "wrong format of id"
	cases[1].err = "request body validation failed ==> wrong UUID provided"

	storage := new(storage.HandlerMock)
	session := getSessionMock(storage)
	put := NewPut(session)

	// do the test
	for _, v := range cases {
		err := put.validateRequestBody(v.body)
		assert.Equal(t, v.err, err.Error(), "Wrong error returned")
	}
}

type dataProviderValidateUnhappy struct {
	body *request.Update
	err string
}

func TestEndpointSessionPutMergeDataHappy(t *testing.T) {
	// setup
	cases := getTestCasesForMergeData()
	storage := new(storage.HandlerMock)
	session := getSessionMock(storage)
	put := NewPut(session)

	// test
	for k, v := range cases {
		res := put.mergeData(v.old, v.new)
		assert.Equal(t, v.resultLength, len(res), fmt.Sprintf("Case %d: Expected length of map is wrong, merge didn't work", k))
		assert.Equal(t, v.result, res, fmt.Sprintf("Case %d: Result is not the expected one",k ))
	}
}

type dataProviderMergeData struct {
	old map[string]string
	new map[string]string
	result map[string]string
	resultLength int
}

func getTestCasesForMergeData() []dataProviderMergeData {
	cases := make([]dataProviderMergeData, 7)

	// case 0: old = empty(nil), new = empty(nil) ==> result = empty
	cases[0].result = map[string]string{}

	// case 1: old = empty, new = empty ==> result = empty
	cases[1].old = map[string]string{}
	cases[1].new = map[string]string{}
	cases[1].result = map[string]string{}

	// case 2: old = empty, new = new data ==> result = new data
	cases[2].old = map[string]string{}
	cases[2].new = map[string]string{"key": "value"}
	cases[2].result = map[string]string{"key": "value"}
	cases[2].resultLength = 1

	// case 3: old = old data, new = empty ==> result = old data
	cases[3].old = map[string]string{"key1": "value1", "key2": "value2"}
	cases[3].new = map[string]string{}
	cases[3].result = map[string]string{"key1": "value1", "key2": "value2"}
	cases[3].resultLength = 2

	// case 4: old = old data, new = new data (different keys) ==> result = old + new data
	cases[4].old = map[string]string{"oldKey": "old value"}
	cases[4].new = map[string]string{"newKey": "new value"}
	cases[4].result = map[string]string{"oldKey": "old value", "newKey": "new value"}
	cases[4].resultLength = 2

	// case 5: old = old data, new = new data (same keys) ==> result = new data
	cases[5].old = map[string]string{"myKey": "old value"}
	cases[5].new = map[string]string{"myKey": "new value"}
	cases[5].result = map[string]string{"myKey": "new value"}
	cases[5].resultLength = 1

	// case 6: old = old data, new = new data (some new keys and some same keys) ==> result = old data (only keys missing in new data) + new data
	cases[6].old = map[string]string{"oldKey": "old value", "sameKey": "another old value"}
	cases[6].new = map[string]string{"newKey": "new value", "sameKey": "another new value"}
	cases[6].result = map[string]string{"oldKey": "old value", "newKey": "new value", "sameKey": "another new value"}
	cases[6].resultLength = 3

	return cases
}

func TestEndpointSessionPutHandlerHappy(t *testing.T) {
	cases := getTestCasesHandlerHappy()

	for k, v := range cases {
		// setup
		storage := v.storage
		session := getSessionMock(storage)
		put := NewPut(session)

		req, err := http.NewRequest("PUT", "/session/", v.body)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(put.Handler)

		// test
		handler.ServeHTTP(rr, req)
		assert.Equal(t, v.resultCode, rr.Code, fmt.Sprintf("Case %d: The response should show a success", k))

		body := new(response.Session)
		decoder := json.NewDecoder(rr.Body)
		err = decoder.Decode(body)
		if err != nil {
			t.Fatal(err)
		}
		if v.id == "" {
			assert.Equal(t, 36, len(body.Id), fmt.Sprintf("Case %d: Id was not returned or wrong id was returned", k))
		} else {
			assert.Equal(t, v.id, body.Id, fmt.Sprintf("Case %d: Id was not returned or wrong id was returned", k))
		}

		data, err := json.Marshal(body.Data)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, v.resultData, data, "Data was not returned or wrong returned")
		assert.Equal(t, 1800, body.Lifetime, "Lifetime was not returned or wrong returned")
		expires := testify.AssertTime{v.expires, body.Expires}
		assert.Condition(t, expires.GreaterThanOrEqual, "Expires needs to be greater or equal than now + default lifetime")
		assert.Equal(t, "", body.Domain, "Domain is not supported yet and should be therefor an empty string")

		storage.AssertExpectations(t)
	}
}

func getTestCasesHandlerHappy() []dataProviderHandler {
	lifetime := time.Duration(1800)
	cases := make([]dataProviderHandler, 2)

	// case 0: with id & body ==> result: same id and a body
	cases[0].id = "8d9af075-1aa6-46c0-913d-ff42f22ca307"
	cases[0].body = strings.NewReader(`{"id": "8d9af075-1aa6-46c0-913d-ff42f22ca307", "data": {"key": "value"}}`)
	cases[0].resultData = []byte(`{"key":"value"}`)
	storage0 := new(storage.HandlerMock)
	storage0.On("Get", cases[0].id).Return("{}", nil)
	storage0.On("Set", cases[0].id, cases[0].resultData, lifetime * time.Second).Return(nil)
	cases[0].storage = storage0
	cases[0].resultCode = http.StatusOK
	cases[0].expires = time.Now().Unix() + response.LIFETIME

	// case 1: with body only ==> result: new id and a body
	cases[1].id = ""
	cases[1].body = strings.NewReader(`{"data":{"myKey":"my value"}}`)
	cases[1].resultData = []byte(`{"myKey":"my value"}`)
	storage1 := new(storage.HandlerMock)
	storage1.On("Set", mock.Anything, cases[1].resultData, lifetime * time.Second).Return(nil)
	cases[1].storage = storage1
	cases[1].resultCode = http.StatusCreated
	cases[1].expires = time.Now().Unix() + response.LIFETIME

	return cases
}

type dataProviderHandler struct {
	storage    *storage.HandlerMock
	id         string
	body       io.Reader
	resultCode int
	resultData []byte
	expires int64
	message string
}


func TestEndpointSessionPutHandlerUnhappy(t *testing.T) {
	cases := getTestCasesHandlerUnhappy()

	for k, v := range cases {
		// setup
		storage := v.storage
		session := getSessionMock(storage)
		put := NewPut(session)

		req, err := http.NewRequest("PUT", "/session/", v.body)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(put.Handler)

		// test
		handler.ServeHTTP(rr, req)
		assert.Equal(t, v.resultCode, rr.Code, fmt.Sprintf("Case %d: The response should show an client or server error", k))
		assert.Equal(t, v.message, rr.Body.String(), fmt.Sprintf("Case %d: Result body should contain an error message", k))

	}
}

func getTestCasesHandlerUnhappy() []dataProviderHandler {
	id := "8d9af075-1aa6-46c0-913d-ff42f22ca307"
	lifetime := time.Duration(1800)
	cases := make([]dataProviderHandler, 4)

	// case 0: getRequestBody fails
	cases[0].body = strings.NewReader("no JSON")
	cases[0].resultCode = http.StatusBadRequest
	cases[0].message = BadRequestText

	// case 1: validateRequestBody fails
	cases[1].body = strings.NewReader(`{"id":"not a valid id"}`)
	cases[1].resultCode = http.StatusBadRequest
	cases[1].message = BadRequestText

	// case 2: loadData fails
	cases[2].body = strings.NewReader(`{"id": "8d9af075-1aa6-46c0-913d-ff42f22ca307", "data": {"key": "value"}}`)
	cases[2].storage = new(storage.HandlerMock)
	cases[2].storage.On("Get", id).Return("", errors.New("Failing loading data"))
	cases[2].resultCode = http.StatusNotFound
	cases[2].message = SessionNotFoundText

	// case 3: storeData fails
	cases[3].body = strings.NewReader(`{"data": {"key": "value"}}`)
	cases[3].storage = new(storage.HandlerMock)
	cases[3].storage.On("Set", mock.Anything, []byte(`{"key":"value"}`), lifetime * time.Second).Return(errors.New("Failing storing data"))
	cases[3].resultCode = http.StatusInternalServerError
	cases[3].message = InternalServerErrorText

	return cases
}

func getSessionMock(storage storage.Handler) *Session {
	auth := &authentication.Authentication{}
	conf := &configuration.Service{}
	return NewSession(storage, auth, conf)
}
