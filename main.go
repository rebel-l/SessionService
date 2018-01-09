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
	// response.PingSummary
	//ps := pingSummaryExample()

	// response.Ping
	//pingExample(ps)


	// do some redis operations if redis flag is set
	//if *runRedis {
	//	fmt.Println("Redis ...")
	//	setEntry("name", "Lars")
	//	setEntry("age", "29")
	//	fmt.Printf("\tHello %s!\n", getEntry("name"))
	//	fmt.Printf("\tYou are %s years old!\n", getEntry("age"))
	//	fmt.Println("")
	//}

	// start to serve
	//if *runServer {
	//	fmt.Println("Static Server ...")
	//	fmt.Println("")
	//}
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
	sessionGet := http.HandlerFunc(sessionGet)	// TODO: remove
	s.router.Handle("/session/", s.middleWare.Middleware(sessionGet)).Methods(http.MethodGet) // TODO: remove

	return s
}

//func pingSummaryExample() *response.PingSummary {
//	fmt.Println("resonse.PingSummary ...")
//	ps := response.NewPingSummary()
//	fmt.Printf("\tStruct (Before): %#v\n", ps)
//	ps.TurnServiceOnline()
//	fmt.Printf("\tService Online: %s\n", ps.Service)
//	ps.TurnStorageOnline()
//	fmt.Printf("\tStorage Online: %s\n", ps.Storage)
//	fmt.Printf("\tStruct (After): %#v\n", ps)
//	res, _ := json.Marshal(ps)
//	fmt.Printf("\tJSON Output %s\n", string(res))
//	fmt.Println("")
//	return ps
//}

//func pingExample(ps *response.PingSummary) {
//	fmt.Println("response.Ping ...")
//	p := response.NewPing()
//	fmt.Printf("\tPing Struct (Before): %#v\n", p)
//	p.Summary = ps
//	p.Notify()
//	fmt.Printf("\tPing Struct (After): %#v\n", p)
//	fmt.Println("")
//}

func (s *server) serve() {
	log.Infof("Listening on port %d ...", s.config.Service.Port)
	err := http.ListenAndServe(":" + strconv.Itoa(s.config.Service.Port), s.router)
	if  err != nil {
		log.Panicf("Couldn't start server. Error: %s", err)
	}
}

func sessionGet(w http.ResponseWriter, r *http.Request) {
	log.Info("Executing session GET")
	w.WriteHeader(http.StatusOK)
	i,err := w.Write([]byte("OK"))
	if i < 1 {
		log.Errorf("Wasn't able to write body: %d", i)
	} else if err != nil {
		log.Errorf("Wasn't able to write body: %s", err)
	}
	log.Info("Executing session GET done!")
}



//func setEntry(key string, value string) {
//	err := getRedisClient().Set(key, value, 0)
//	if err.Err() != nil {
//		panic(err)
//	}
//}

//func getEntry(key string) string {
//	result, err := getRedisClient().Get(key).Result()
//	if err != nil {
//		panic(err)
//	}
//
//	return result
//}