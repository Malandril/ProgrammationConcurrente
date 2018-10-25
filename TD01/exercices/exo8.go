package main

import (
	"fmt"
)

func main() {
	fmt.Println("Prime numbers")
	out := primeNumbers()
	for n := <-out; n < 100; n = <-out { // Affichage des nombres premiers < 100
		fmt.Println(n)
	}
	fmt.Println(<-out)
}

func primeNumbers() chan int {
	out := make(chan int)
	go func() {
		numbers := generateNumbers()
		crible(numbers, out)
	}()
	return out
}

func crible(numbers chan int, out chan int) {
	p := <-numbers
	out <- p
	primes := make(chan int)
	go func() {
		for x := range numbers {
			if x%p != 0 {
				primes <- x
			}
		}
		close(primes)
	}()
	crible(primes, out)
}
func generateNumbers() chan int {
	c := make(chan int)
	go func() {
		for i := 2; ; i++ {
			c <- i
		}
		close(c)
	}()
	return c
}
