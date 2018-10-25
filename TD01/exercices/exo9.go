package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

var PORT = "1234"

const timeout = 10 * time.Second

type Client struct {
	name string
	conn net.Conn
}
type Message struct {
	sender  Client
	message string
}

func sendDisconnection(client Client, disconnection chan Client) {
	fmt.Println("AAAAAAAAAAAAAAAAAAAAAAAAAAAA")
	disconnection <- client
}

// Handles a client
func handleConnection(conn net.Conn, messageChan chan Message, connection chan Client, disconnection chan Client) {
	reader := bufio.NewReader(conn)
	fmt.Fprintf(conn, "Who are you?: ")
	login, err := reader.ReadString('\n')
	if err != nil {
		log.Println(err)
		conn.Close()
		return
	}
	login = strings.TrimSuffix(login, "\n")
	client := Client{login, conn}
	connection <- client
	defer sendDisconnection(client, disconnection)

	// channel to allow the select to know when the client did not send messages since a long time
	receivedMessage := make(chan bool)

	// reads data from the client
	go func() {
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				log.Println(err)
				receivedMessage <- false
				return
			}
			receivedMessage <- true
			messageChan <- Message{client, line}
		}
	}()

	// If the server does not receive data from the client disconnect it with defer
	for {
		select {
		case x := <-receivedMessage:
			if !x {
				return
			}
		case <-time.After(timeout):
			return
		}
	}

}

// This go routine sends the messages to all clients
// The client map is only modified here in the select so there is no concurrent modification or access
func messageSender(clients map[string]net.Conn, messageChan chan Message, connection chan Client, disconnection chan Client) {
	for {
		select {
		case message := <-messageChan:
			fmt.Print(message.sender.name, "sent", message.message)
			for _, conn := range clients {
				fmt.Fprintf(conn, "%s: %s", message.sender.name, message.message)
			}
		case client := <-connection:
			_, ok := clients[client.name]
			if ok {
				fmt.Fprintln(client.conn, "Login already used")
				return
			}
			fmt.Println(client.name, "connected")
			clients[client.name] = client.conn
			for _, conn := range clients {
				fmt.Fprintf(conn, "%s connected\n", client.name)
			}
		case client := <-disconnection:
			client.conn.Close()
			delete(clients, client.name)
			fmt.Println(client.name, "disconnected")
			for _, conn := range clients {
				fmt.Println(client.name)
				fmt.Fprintf(conn, "%s disconnected\n", client.name)
			}
		}

	}
}

func mainLoop() {
	listener, err := net.Listen("tcp", "localhost:"+PORT)
	if err != nil {
		log.Fatal(err)
	}
	//Message channel
	messageChan := make(chan Message)
	//New connection channel
	connection := make(chan Client)
	// Disconnection channel
	disconnection := make(chan Client)

	//Client map to store all clients
	clients := make(map[string]net.Conn)

	go messageSender(clients, messageChan, connection, disconnection)
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go handleConnection(conn, messageChan, connection, disconnection)
	}
}

func main() {
	mainLoop()
}
