package main

import (
	"fmt"
	"time"
)

func main() {
	c := make(chan int)
	finished := make(chan bool)
	go affichage(c, finished)

	for i := 0; i < 10; i++ {
		c <- i
	}
	close(c)
	<-finished
}

func affichage(c chan int, finished chan bool) {
	for x := range c {
		time.Sleep(20 * time.Millisecond)
		fmt.Println(x)
	}
	finished <- true
}
