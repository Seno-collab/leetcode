package main

import (
	"fmt"
	"time"
)

func sayHello() {
	fmt.Println("Hello from goroutine")
}

func main() {
	go sayHello()
	time.Sleep(200 * time.Second)
	fmt.Println("Hello from main")
}
