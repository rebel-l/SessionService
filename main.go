package main

import (
	"fmt"
	//"github.com/go-redis/redis"
	"github.com/rebel-l/sessionservice/src/configuration"
	"github.com/rebel-l/sessionservice/src/endpoint"
	//"github.com/rebel-l/sessionservice/src/response"
	log "github.com/sirupsen/logrus"
	"net/http"
	//"encoding/json"
	"strconv"

)

func main() {
	fmt.Println("")

	// response.PingSummary
	//ps := pingSummaryExample()

	// response.Ping
	//pingExample(ps)

	// init config
	//fmt.Println("Config init ...")
	config := configuration.Init()
	initLogging(config.Service.LogLevel)
	//fmt.Printf("\tConfig.Service: %#v\n", config.Service)
	//fmt.Printf("\tConfig.Redis: %#v\n", config.Redis)
	//fmt.Println("")

	// do some redis operations if redis flag is set
	//if *runRedis {
	//	fmt.Println("Redis ...")
	//	redisPing()
	//	setEntry("name", "Lars")
	//	setEntry("age", "29")
	//	fmt.Printf("\tHello %s!\n", getEntry("name"))
	//	fmt.Printf("\tYou are %s years old!\n", getEntry("age"))
	//	fmt.Println("")
	//}

	// start to serve
	//if *runServer {
	//	fmt.Println("Static Server ...")
		serve(*config.Service)
	//	fmt.Println("")
	//}
}

func initLogging(loglevel log.Level) {
	log.SetLevel(loglevel)
	log.Debug("Logging initialized")
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

func serve(config configuration.Service) {
	endpoint.InitDocsEndpoint()
	endpoint.InitPing()
	log.Infof("Listening on port %d ...", config.Port)
	err := http.ListenAndServe(":" + strconv.Itoa(config.Port), nil)
	if  err != nil {
		log.Panicf("Couldn't start server. Error: %s", err)
	}
}

//func getRedisClient() *redis.Client {
//	client := redis.NewClient(&redis.Options{
//		Addr: "redis:6379",
//		Password: "",
//		DB: 0,
//	})
//
//	return client
//}

//func redisPing() {
//	pong, err := getRedisClient().Ping().Result()
//	fmt.Println(pong, err)
//}

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