package endpoint

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

type ResponseWriterMock struct {
	HeaderMock      *HeaderMock
	WriteHeaderMock *WriteHeaderMock
	WriteMock       *WriteMock
}

type HeaderMock struct {
	counter int
	header *http.Header
}

type WriteHeaderMock struct {
	counter int
	code int
}

type WriteMock struct {
	counter int
	data string
	err error
	dataLen int
}

type PlainTestCase struct {
	code int
	msg string
	mock *ResponseWriterMock
}

func dataProviderPlainTestCases() []PlainTestCase {
	dataProvider := make([]PlainTestCase, 3)

	// TestCase 1. Happy Path
	dataProvider[0].code = 123
	dataProvider[0].msg = "MyMessage"
	wm1 := getWriterMock(len(dataProvider[0].msg), nil)
	dataProvider[0].mock= getResponseWriterMock(wm1)

	// TestCase 2. Error Unhappy Path
	dataProvider[1].code = 333
	dataProvider[1].msg = "Error"
	wm2 := getWriterMock(len(dataProvider[1].msg), errors.New("Something went wrong"))
	dataProvider[1].mock = getResponseWriterMock(wm2)

	// TestCase 3. Data Length 0 Unhappy Path
	dataProvider[2].code = 444
	dataProvider[2].msg = "Data Length 0"
	wm3 := getWriterMock(0, nil)
	dataProvider[2].mock = getResponseWriterMock(wm3)

	return dataProvider
}

func (rwm ResponseWriterMock) Write(data []byte) (int, error) {
	rwm.WriteMock.counter++
	rwm.WriteMock.data = string(data)
	return rwm.WriteMock.dataLen, rwm.WriteMock.err
}

func (rwm ResponseWriterMock) Header() http.Header {
	rwm.HeaderMock.counter++
	header := http.Header{}
	rwm.HeaderMock.header = &header
	return header
}

func (rwm ResponseWriterMock) WriteHeader(code int) {
	rwm.WriteHeaderMock.counter++
	rwm.WriteHeaderMock.code = code
}

func TestMessageNew(t *testing.T) {
	fixture := ResponseWriterMock{}
	msg := NewMessage(fixture)
	assert.Equal(t, fixture, msg.res, "Message has not the right response writer set")
}

func getResponseWriterMock(wm *WriteMock) *ResponseWriterMock {
	mock := new(ResponseWriterMock)
	mock.HeaderMock = new(HeaderMock)
	mock.WriteHeaderMock = new(WriteHeaderMock)
	mock.WriteMock = wm
	return mock
}

func getWriterMock(length int, err error) *WriteMock {
	wm := new(WriteMock)
	wm.dataLen = length
	wm.err = err
	return wm
}

func TestMessagePlainHappy(t *testing.T) {
	for _, v := range dataProviderPlainTestCases() {
		msg := NewMessage(v.mock)
		msg.Plain(v.msg, v.code)

		// Header
		testErrMsg := fmt.Sprintf("Message has not set header for '%s: %s'", ContentHeader, ContentTypePlain)
		assert.Equal(t, 1, v.mock.HeaderMock.counter, "Header not set")
		assert.Equal(t, ContentTypePlain, v.mock.HeaderMock.header.Get(ContentHeader), testErrMsg)

		// WriteHeader
		assert.Equal(t, 1, v.mock.WriteHeaderMock.counter, "Header not written")
		assert.Equal(t, v.code, v.mock.WriteHeaderMock.code, "Header code not written")

		// Write
		assert.Equal(t, 1, v.mock.WriteMock.counter, "Data not written")
		assert.Equal(t, v.msg, v.mock.WriteMock.data, "Data not written")
	}
}

func TestInternalServerErrorHappy(t *testing.T) {
	fixture := "Internal Server Error"
	wm := getWriterMock(len(fixture), nil)
	mock := getResponseWriterMock(wm)
	msg := NewMessage(mock)
	msg.InternalServerError(fixture)

	// Header
	testErrMsg := fmt.Sprintf("Message has not set header for '%s: %s'", ContentHeader, ContentTypePlain)
	assert.Equal(t, 1, mock.HeaderMock.counter, "Header not set")
	assert.Equal(t, ContentTypePlain, mock.HeaderMock.header.Get(ContentHeader), testErrMsg)

	// WriteHeader
	assert.Equal(t, 1, mock.WriteHeaderMock.counter, "Header not written")
	assert.Equal(t, http.StatusInternalServerError, mock.WriteHeaderMock.code, "Header code not written")

	// Write
	assert.Equal(t, 1, mock.WriteMock.counter, "Data not written")
	assert.Equal(t, fixture, mock.WriteMock.data, "Data not written")
}

func TestSendJson(t *testing.T) {
	data := make(map[string]string)
	data["key1"] = "value 1"
	data["key2"] = "value 2"

	fixture := "{\"key1\":\"value 1\",\"key2\":\"value 2\"}\n"

	wm := getWriterMock(len(fixture), nil)
	mock := getResponseWriterMock(wm)
	msg := NewMessage(mock)
	msg.SendJson(data, 200)

	// Header
	testErrMsg := fmt.Sprintf("Message has not set header for '%s: %s'", ContentHeader, ContentTypePlain)
	assert.Equal(t, 1, mock.HeaderMock.counter, "Header not set")
	assert.Equal(t, ContentTypeJson, mock.HeaderMock.header.Get(ContentHeader), testErrMsg)

	// WriteHeader
	assert.Equal(t, 1, mock.WriteHeaderMock.counter, "Header not written")
	assert.Equal(t, 200, mock.WriteHeaderMock.code, "Header code not written")

	// Write
	assert.Equal(t, 1, mock.WriteMock.counter, "Data not written")
	assert.Equal(t, fixture, mock.WriteMock.data, "Data not written")
}
