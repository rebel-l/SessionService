package endpoint

import (
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

func TestMessagePlainHappy(t *testing.T) {
	fixture := "MyMessage"
	mock := ResponseWriterMock{}
	mock.HeaderMock = new(HeaderMock)
	mock.WriteHeaderMock = new(WriteHeaderMock)
	wm := new(WriteMock)
	wm.dataLen = len(fixture)
	wm.err = nil
	mock.WriteMock = wm
	msg := NewMessage(mock)
	msg.Plain(fixture, 123)

	// Header
	testErrMsg := fmt.Sprintf("Message has not set header for '%s: %s'", ContentHeader, ContentTypePlain)
	assert.Equal(t, 1, mock.HeaderMock.counter, "Header not set")
	assert.Equal(t, ContentTypePlain, mock.HeaderMock.header.Get(ContentHeader), testErrMsg)

	// WriteHeader
	assert.Equal(t, 1, mock.WriteHeaderMock.counter, "Header not written")
	assert.Equal(t, 123, mock.WriteHeaderMock.code, "Header code not written")

	// Write
	assert.Equal(t, 1, mock.WriteMock.counter, "Data not written")
	assert.Equal(t, fixture, mock.WriteMock.data, "Data not written")
}

// TODO: write unhappy 1. writeMock returns error & 2. writeMock returns len < 1
