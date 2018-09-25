package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("salut")
	Tick(3 * time.Second)
	fmt.Println("sleepy")
	sleepy(3 * time.Second)
	fmt.Println("tick")

}

func Tick(t time.Duration) {
	c := make(chan bool)
	go func() {
		for {
			fmt.Print(".")
			select {
			case <-c:
				return
			case <-time.After(1 * time.Second):
			}
		}
	}()
	<-time.After(t)
	close(c)
}

func sleepy(t time.Duration) {
	<-time.After(t)
}
