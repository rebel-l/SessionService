package main

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
	"github.com/rebel-l/sessionservice/src/authentication"
	"github.com/rebel-l/sessionservice/src/configuration"
	"github.com/rebel-l/sessionservice/src/endpoint"
	"github.com/rebel-l/sessionservice/src/endpoint/session"
	"github.com/rebel-l/sessionservice/src/storage"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

func main() {
	fmt.Println("")

	// init config & server
	parser := configuration.GetParser()
	server := new(server)
	server.config = parser.Parse()
	server.init().serve()
}

type server struct {
	config *configuration.Config
	router *mux.Router
	middleWare *authentication.Authentication
	redis *redis.Client	// ToDo: deprecated
	storage storage.Handler
}

func (s *server) init() *server {
	// init logging
	log.SetLevel(s.config.Service.LogLevel)
	log.Debugf("Logging initialized: %d", s.config.Service.LogLevel)

	// init Router
	s.router = mux.NewRouter()

	// init middleware
	s.middleWare = authentication.New(s.config.AccountList)

	// init storage & endpoints
	s.initStorage().
		initEndpoints()

	return s
}

func (s *server) initStorage() *server {
	s.redis = redis.NewClient(s.config.Redis)
	s.storage = storage.NewRedis(s.redis)
	return s
}

func (s *server) initEndpoints() *server {
	// docs
	endpoint.InitDocsEndpoint(s.router)

	// ping
	endpoint.InitPing(s.redis, s.router)

	// session
	sessionEndpoint := session.NewSession(s.storage, s.middleWare, s.config.Service)
	sessionEndpoint.Init(s.router)

	return s
}

func (s *server) serve() {
	log.Infof("Listening on port %d ...", s.config.Service.Port)
	err := http.ListenAndServe(":" + strconv.Itoa(s.config.Service.Port), s.router)
	if  err != nil {
		log.Panicf("Couldn't start server. Error: %s", err)
	}
}
