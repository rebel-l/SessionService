package main

import "fmt"

func main(){
	var a = multiply(6,7)
	fmt.Printf("Hello World\n")
	fmt.Printf("6 x 7 is %d\n", a)
}

func multiply(a int, b int) int {
	return a * b
}
