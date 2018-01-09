package session

import (
	"testing"
	"github.com/go-redis/redis"
	"github.com/rebel-l/sessionservice/src/authentication"
	"github.com/rebel-l/sessionservice/src/configuration"
	"github.com/stretchr/testify/assert"
	"github.com/gorilla/mux"
	"net/http"
	"reflect"
	"unsafe"
	"fmt"
)

func TestEndpointSessionNewSession(t *testing.T) {
	client := redis.NewClient(&redis.Options{})
	auth := &authentication.Authentication{}
	conf := &configuration.Service{}
	session := NewSession(client, auth, conf)
	assert.Equal(t, client, session.Redis, "Wasn't able to set Redis client")
	assert.Equal(t, auth, session.Authentication, "Wasn't able to set authentication")
	assert.Equal(t, conf, session.Config, "Wasn't able to set config")
}

func TestEndpointSessionHandlerFactoryHappy(t *testing.T) {
	client := redis.NewClient(&redis.Options{})
	auth := &authentication.Authentication{}
	conf := &configuration.Service{}
	session := NewSession(client, auth, conf)
	handler := session.handlerFactory(http.MethodPut)
	result := reflect.TypeOf(handler).String()
	assert.Equal(t, "http.HandlerFunc", result, "Returned type needs to be of type 'http.HandlerFunc'")
}

func TestEndpointSessionInit(t *testing.T) {
	expectedMethods := make(map[string]int)
	expectedMethods["PUT"] = 0

	client := redis.NewClient(&redis.Options{})
	auth := &authentication.Authentication{}
	conf := &configuration.Service{}
	session := NewSession(client, auth, conf)
	router := mux.NewRouter()

	// assert before initialisation
	routes := extractRoutes(router)
	assert.Equal(t, 0, len(routes), "There should not exist any route before initialisation")

	// assert after initialisation
	session.Init(router)
	routes = extractRoutes(router)
	assert.Equal(t, 1, len(routes), "There should exactly one route exist after initialisation")
	for _, v := range routes {
		methods, err := v.GetMethods()
		if err != nil {
			t.Fatal("Unrecoverable error on getting routers methods")
		}
		for _, method := range methods {
			expectedMethods[method]++
		}
	}

	for k, v := range expectedMethods {
		assert.Exactly(t, 1, v, fmt.Sprintf("%s method shoul only be set up once", k))
	}
}

func extractRoutes(router *mux.Router) []*mux.Route {
	r := reflect.ValueOf(router).Elem()
	f := r.FieldByName("routes")

	// cheat to access unexported values
	f = reflect.NewAt(f.Type(), unsafe.Pointer(f.UnsafeAddr())).Elem()

	// convert routes
	return f.Interface().([]*mux.Route)
}
