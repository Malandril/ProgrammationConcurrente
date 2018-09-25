package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"time"
)

func main() {
	runtime.GOMAXPROCS(1)
	now := time.Now()
	c := make(chan int)
	number := 1000
	for i := 0; i < number; i++ {
		go gogadgeto(c)
	}
	for i := 0; i < number; i++ {
		<-c
	}
	fmt.Printf("%s\n", time.Since(now))
}

func gogadgeto(c chan int) {
	wait := rand.Intn(500)
	time.Sleep(time.Duration(wait) * time.Millisecond)
	fmt.Println("salut")
	c <- 1
}
