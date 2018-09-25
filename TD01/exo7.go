package main

import (
	"fmt"
	"strconv"
	"sync"
)

func main() {
	g1 := generator(10)
	g2 := generator(5)
	g3 := generator(7)
	strings := merge(g1, g2, g3)
	for x := range strings {
		fmt.Printf("p %s \n",x)
	}
}
func generator(max int) chan string {
	c := make(chan string)
	go func() {
		for i := 0;i<max ; i++ {
			c <- "test: " + strconv.Itoa(i)
		}
		close(c)
	}()
	return c
}

func merge(chans ...chan string) chan string {
	res := make(chan string)
	var wg sync.WaitGroup
	for _, c := range chans {
		go func(c chan string) {
			wg.Add(1)
			for s := range c {
				res <- s
			}
			wg.Done()
		}(c)
	}
	go func() {
		wg.Wait()
		close(res)
	}()
	return res
}
