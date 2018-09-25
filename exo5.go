package main

import "fmt"

func main() {
	baton := make(chan int)
	quit := make(chan int)

	// Faire partir le premier coureur
	go coureur(baton, quit)
	baton <- 1
	<-quit
	fmt.Println("Fin de la course")
}

func coureur(baton chan int, quit chan int) {
	b := <-baton
	fmt.Printf("baton %d", b)
	quit<-2
}