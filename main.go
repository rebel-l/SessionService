package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/rebel-l/sessionservice/src/authentication"
	"github.com/rebel-l/sessionservice/src/configuration"
	"github.com/rebel-l/sessionservice/src/endpoint"
	"github.com/rebel-l/sessionservice/src/request"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"github.com/rebel-l/sessionservice/src/response"
	"github.com/gorilla/mux"
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
	middleWare *authentication.Authentification
	redis *redis.Client
}

func (s *server) init() *server {
	// init logging
	log.SetLevel(s.config.Service.LogLevel)
	log.Debugf("Logging initialized: %d", s.config.Service.LogLevel)

	// init Router
	s.router = mux.NewRouter()

	// init middleware
	s.middleWare = authentication.New(s.config.AccountList)

	// init storage
	s.redis = redis.NewClient(s.config.Redis)

	// init endpoints
	s.initEndpoints()

	return s
}

func (s *server) initEndpoints() *server {
	// docs
	endpoint.InitDocsEndpoint(s.router)

	// ping
	endpoint.InitPing(s.redis, s.router)

	// session
	sessionGet := http.HandlerFunc(sessionGet)
	sessionPut := http.HandlerFunc(sessionPut)
	s.router.Handle("/session/", s.middleWare.Middleware(sessionGet)).Methods(http.MethodGet)
	s.router.Handle("/session/", s.middleWare.Middleware(sessionPut)).Methods(http.MethodPut)
	log.Debug("Session endpoint initialized")

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
	i,_ := w.Write([]byte("OK"))
	if i < 1 {
		log.Errorf("Wasn't able to write body: %d", i)
	}
	log.Info("Executing session GET done!")
}

func sessionPut(w http.ResponseWriter, r *http.Request) {
	log.Println("Executing session PUT")

	// read request body
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	var data request.Update
	err := decoder.Decode(&data)
	if err != nil {
		log.Errorf("Unable to read request body: %s", err)
	}

	log.Debugf("Id to update: %s", data.Id)
	for key, value := range data.Data {
		log.Debugf("%s: %s", key, value)
	}

	// write request
	w.WriteHeader(http.StatusOK)
	session := response.NewSession("",0)
	err = json.NewEncoder(w).Encode(session)
	if err != nil {
		log.Errorf("Wasn't able to write body: %s", err)
	}

	log.Info("Executing session PUT done!")
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