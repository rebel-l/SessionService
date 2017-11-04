package authentication

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestMiddlewareDefaultValue(t *testing.T) {
	allowedAccounts := make(map[string]Account)
	allowedAccounts["MyApp"] = Account{"token"}
	m := New(allowedAccounts)
	assert.Equal(t, "token", m.allowedAccounts["MyApp"].ApiKey, "AllowedAccounts was not set")
}

func TestMiddlewareAuthenticateHappy(t *testing.T) {
	allowedAccounts := make(map[string]Account)
	allowedAccounts["App"] = Account{"tokenizer"}
	m := New(allowedAccounts)
	assert.True(t, m.authenticate("App", "tokenizer"), "Authentication should be successful with correct AppId and ApiKey")
}

func TestMiddlewareAuthenticateUnhappy(t *testing.T) {
	for _, testCase := range MiddlewareAuthenticateUnhappyDataProvider() {
		assert.False(t, testCase.auth.authenticate(testCase.appId, testCase.apiKey), "Authentication should fail!")
	}
}

type AuthenticateUnhappyTestCases struct {
	auth *Authentification
	appId string
	apiKey string
}

func MiddlewareAuthenticateUnhappyDataProvider() []AuthenticateUnhappyTestCases {
	var testCases []AuthenticateUnhappyTestCases

	// 1. empty allowedAccounts
	testCase1 := AuthenticateUnhappyTestCases{
		auth: new(Authentification),
		appId: "noAccounts",
		apiKey: "yes1",
	}
	testCases = append(testCases, testCase1)

	// 2. missing AppId in allowedAccounts
	allowedAccounts2 := make(map[string]Account)
	allowedAccounts2["existingApp"] = Account{"yes2"}
	testCase2 := AuthenticateUnhappyTestCases{
		auth: New(allowedAccounts2),
		appId: "missingApp",
		apiKey: "yes2",
	}
	testCases = append(testCases, testCase2)

	// 3. empty AppId
	allowedAccounts3 := make(map[string]Account)
	allowedAccounts3["emptyApp"] = Account{"yes3"}
	testCase3 := AuthenticateUnhappyTestCases{
		auth: New(allowedAccounts3),
		appId: "",
		apiKey: "yes3",
	}
	testCases = append(testCases, testCase3)

	//4. right AppId, wrong ApiKey
	allowedAccounts4 := make(map[string]Account)
	allowedAccounts4["AppId"] = Account{"yes4"}
	testCase4 := AuthenticateUnhappyTestCases{
		auth: New(allowedAccounts4),
		appId: "AppId",
		apiKey: "no",
	}
	testCases = append(testCases, testCase4)

	return testCases
}

func TestMiddlewareHandlerHappy(t *testing.T) {
	allowedAccounts := make(map[string]Account)
	allowedAccounts["NiceApp"] = Account{"handleIt"}
	m := New(allowedAccounts)
	httpHandlerMock := new(MiddlewareHttpHandlerMock)
	h := m.Middleware(httpHandlerMock)
	req, err := http.NewRequest("GET", "/session/", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set(HEADER_APPID_KEY, "NiceApp")
	req.Header.Set(HEADER_APIKEY_KEY, "handleIt")
	wm := new(MiddlewareHttpResponseWriterMock)
	h.ServeHTTP(wm, req)
	assert.Equal(t, 1, httpHandlerMock.nextHandlercounter, "Successful authentication needs to call next handler")
	assert.Equal(t, 0, wm.writeHeaderCounter, "Successful authentication should not manipulate the response header")
	assert.Equal(t, 0, wm.writeCounter, "Successful authentication should not manipulate the response body")
}

func TestMiddlewareHandlerUnhappy(t *testing.T) {
	allowedAccounts := make(map[string]Account)
	allowedAccounts["CoolApp"] = Account{"wrongKey"}
	m := New(allowedAccounts)
	httpHandlerMock := new(MiddlewareHttpHandlerMock)
	h := m.Middleware(httpHandlerMock)
	req, err := http.NewRequest("GET", "/session/", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set(HEADER_APPID_KEY, "CoolApp")
	req.Header.Set(HEADER_APIKEY_KEY, "keyDoesNotMatch")
	wm := new(MiddlewareHttpResponseWriterMock)
	h.ServeHTTP(wm, req)
	assert.Equal(t, 0, httpHandlerMock.nextHandlercounter, "Failed authentication should stop handler chain")
	assert.Equal(t, 1, wm.writeHeaderCounter, "Failed authentication should manipulate the response header")
	assert.Equal(t, http.StatusForbidden, wm.writeHeaderValue, "Failed authentification should send 403 http stauts code")
	assert.Equal(t, 1, wm.writeCounter, "Failed authentication should manipulate the response body")
	expectedMsg := "Authetification failed! Ensure you send correct X-APP-ID & X-API-KEY within your header."
	assert.Equal(t, expectedMsg, string(wm.writeValue), "Failed authentication should send a message in response")
}

type MiddlewareHttpHandlerMock struct {
	nextHandlercounter int

}

func (m *MiddlewareHttpHandlerMock) ServeHTTP(res http.ResponseWriter, req *http.Request)  {
	m.nextHandlercounter++
}

type MiddlewareHttpResponseWriterMock struct {
	writeCounter int
	writeValue []byte
	writeHeaderCounter int
	writeHeaderValue int
}

func (m *MiddlewareHttpResponseWriterMock) Header() http.Header {
	return make(http.Header)
}

func (m *MiddlewareHttpResponseWriterMock) Write(msg []byte) (int, error) {
	m.writeValue = msg
	m.writeCounter++
	return 1, nil
}

func (m *MiddlewareHttpResponseWriterMock) WriteHeader(status int) {
	m.writeHeaderValue = status
	m.writeHeaderCounter++
}