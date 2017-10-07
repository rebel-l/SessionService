package main

import (
	"flag"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/rebel-l/sessionservice/src/response"
	"log"
	"net/http"
	"encoding/json"
)

func main() {
	fmt.Println("")

	// response.PingSummary
	ps := pingSummaryExample()

	// response.Ping
	pingExample(ps)

	// parse the flags
	fmt.Println("Parse flags ...")
	runRedis := flag.Bool("redis", false, "execute redis operations")
	runServer := flag.Bool("server", false, "execute server listening")
	flag.Parse()
	fmt.Printf("\tRun Redis: %t\n", *runRedis)
	fmt.Printf("\tRun Server: %t\n", *runServer)
	fmt.Println("")

	// print some texts
	examples()

	// do some redis operations if redis flag is set
	if *runRedis {
		fmt.Println("Redis ...")
		redisPing()
		setEntry("name", "Lars")
		setEntry("age", "29")
		fmt.Printf("\tHello %s!\n", getEntry("name"))
		fmt.Printf("\tYou are %s years old!\n", getEntry("age"))
		fmt.Println("")
	}

	// start to serve
	if *runServer {
		fmt.Println("Static Server ...")
		serveStatic()
		fmt.Println("")
	}
}

func examples() {
	fmt.Println("Print some examples ...")
	a := multiply(6,7)
	fmt.Printf("\tHello World\n")
	fmt.Printf("\t6 x 7 is %d\n", a)
	fmt.Println("")
}

func pingSummaryExample() *response.PingSummary {
	fmt.Println("resonse.PingSummary ...")
	ps := response.NewPingSummary()
	fmt.Printf("\tStruct (Before): %#v\n", ps)
	ps.TurnServiceOnline()
	fmt.Printf("\tService Online: %s\n", ps.Service)
	ps.TurnStorageOnline()
	fmt.Printf("\tStorage Online: %s\n", ps.Storage)
	fmt.Printf("\tStruct (After): %#v\n", ps)
	res, _ := json.Marshal(ps)
	fmt.Printf("\tJSON Output %s\n", string(res))
	fmt.Println("")
	return ps
}

func pingExample(ps *response.PingSummary) {
	fmt.Println("response.Ping ...")
	p := response.NewPing()
	fmt.Printf("\tPing Struct (Before): %#v\n", p)
	p.Summary = ps
	p.Notify()
	fmt.Printf("\tPing Struct (After): %#v\n", p)
	fmt.Println("")
}

func serveStatic() {
	fs := http.FileServer(http.Dir("docs"))
	http.Handle("/docs/", http.StripPrefix("/docs/", fs))
	log.Printf("\tListening ...")
	err := http.ListenAndServe(":4000", nil)
	if  err != nil {
		panic(err)
	}
}

func multiply(a int, b int) int {
	return a * b
}

func getRedisClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr: "redis:6379",
		Password: "",
		DB: 0,
	})

	return client
}

func redisPing() {
	pong, err := getRedisClient().Ping().Result()
	fmt.Println(pong, err)
}

func setEntry(key string, value string) {
	err := getRedisClient().Set(key, value, 0)
	if err.Err() != nil {
		panic(err)
	}
}

func getEntry(key string) string {
	result, err := getRedisClient().Get(key).Result()
	if err != nil {
		panic(err)
	}

	return result
}