package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	fmt.Println("Salut")
	reader := bufio.NewReader(os.Stdin)
	conn, err := net.Dial("tcp", "localhost:1234")
	checkError(err)
	go printIncoming(conn)()
	for {
		str, err := reader.ReadString('\n')
		checkError(err)
		conn.Write([]byte(str))
	}
}

func printIncoming(conn net.Conn) func() {
	return func() {
		reader := bufio.NewReader(conn)
		for {
			line, err := reader.ReadString('\n')
			checkError(err)
			fmt.Printf("%s", line)
		}
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
