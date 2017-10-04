package main

import (
	"fmt"
	"github.com/go-redis/redis"
)

func main(){
	var a = multiply(6,7)
	fmt.Printf("Hello World\n")
	fmt.Printf("6 x 7 is %d\n", a)
	redisPing()
	setEntry("name", "Lars")
	setEntry("age", "29")
	fmt.Printf("Hello %s!\n", getEntry("name"))
	fmt.Printf("You are %s years old!\n", getEntry("age"))
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