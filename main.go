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
