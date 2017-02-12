package main

import (
	"gopkg.in/redis.v5"
	//"fmt"
)

func main()  {
	println("Hello");
	var test  = redis.Nil;
	if (test != redis.Nil) {
		println("Toll");
	}
}