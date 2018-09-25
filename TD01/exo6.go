package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func Producer(out chan int) {
	for i := 0; i < 100; i++ {
		out <- i
	}
	close(out)
}

func Consumer(in chan []int) {
	for result := range in {
		fmt.Printf("Job %02d took %3dms\n", result[0], result[1])
	}
}

func Worker(in chan int, out chan []int, wg sync.WaitGroup) {
	wg.Add(1)
	for n := range in {
		delay := rand.Int() % 100 // Take an arbitrary delay
		time.Sleep(time.Duration(delay) * time.Millisecond)

		out <- []int{n, delay} // result  is [n, delay]
	}
	wg.Done()
}

func main() {
	start := time.Now()

	in := make(chan int)    // Channel of int
	out := make(chan []int) // Channel of couples [n, time(n)]
	var wg sync.WaitGroup
	go Producer(in) // Launch a producer
	for i := 0; i < 10; i++ { // Launch a 10 workers
		go Worker(in, out, wg)
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	Consumer(out) // and, ... consume
	fmt.Println("Total execution time: ", time.Since(start))
}
