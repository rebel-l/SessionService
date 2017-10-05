package main

import (
	"flag"
	"fmt"
	"github.com/go-redis/redis"
	"log"
	"net/http"
)

func main() {
	// parse the flags
	runRedis := flag.Bool("redis", false, "execute redis operations")
	runServer := flag.Bool("server", false, "execute server listening")
	flag.Parse()
	fmt.Printf("Run Redis: %t\n", *runRedis)
	fmt.Printf("Run Server: %t\n", *runServer)

	// print some texts
	a := multiply(6,7)
	fmt.Printf("Hello World\n")
	fmt.Printf("6 x 7 is %d\n", a)

	// do some redis operations if redis flag is set
	if *runRedis {
		redisPing()
		setEntry("name", "Lars")
		setEntry("age", "29")
		fmt.Printf("Hello %s!\n", getEntry("name"))
		fmt.Printf("You are %s years old!\n", getEntry("age"))
	}

	// start to serve
	if *runServer {
		serveStatic()
	}
}

func serveStatic() {
	fs := http.FileServer(http.Dir("docs"))
	http.Handle("/docs/", http.StripPrefix("/docs/", fs))
	log.Printf("Listening ...")
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