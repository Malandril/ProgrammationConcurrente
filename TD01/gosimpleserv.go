//
// gosimpleserv.go  -- A simple server
//
//           Author: Erick Gallesio [eg@unice.fr]
//    Creation date: 17-Oct-2016 14:42 (eg)
// Last file update: 17-Oct-2016 15:50 (eg)
//
// +++
package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

// +slide
func handleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Println(err)
			return
		}
		line = strings.TrimSuffix(line, "\n")
		fmt.Fprintf(conn, "You said => %s\n", line)
		fmt.Printf("Received %s from %s\n", line, conn.RemoteAddr())
	}
}

func main() {
	listener, err := net.Listen("tcp", "localhost:1234")
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go handleConnection(conn)
	}
}
