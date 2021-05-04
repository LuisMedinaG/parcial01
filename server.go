package main

import (
	"log"
	"net"
	"io"
	"encoding/gob"
)

const (
	CONN_HOST = "localhost"
	CONN_PORT = "3333"
	CONN_TYPE = "tcp"
)

var connections = make(map[int]net.Conn) 
var id = 0

func main() {
	l, err := net.Listen(CONN_TYPE, CONN_HOST + ":" + CONN_PORT)
	if err != nil {
		log.Println("Error listening:", err.Error())
		return
	}
	
	defer l.Close()
	
	log.Println("Listening on " + CONN_HOST + ":" + CONN_PORT)
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Println("Error accepting: ", err.Error())
			return
		}
		id++
		connections[id] = conn
		go handleRequest(id, conn)
	}
}

func handleRequest(id int, conn net.Conn) {
	defer func(){
		conn.Close()
		delete(connections, id)
	}()

	for {
		var msg string
		err := gob.NewDecoder(conn).Decode(&msg)
		if err == io.EOF {
			// Close conn and exit
			conn.Close()
			log.Println("Connection Closed.")
			return
		}
		if err != nil {
			log.Println(err)
			return
		}
		// log.Printf("Message Received: %s\n", msg)
		broadcast(conn, msg)
	}
}


func broadcast(conn net.Conn, msg string) {
	for _, c := range connections {
		if c != conn {
		err := gob.NewEncoder(c).Encode(&msg)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

// // ----------------------------------------------------------------

package main

import (
	"fmt"
	"bufio"
	"net"
)

type Client struct {
	// id 		 id int
	incoming chan string
	outgoing chan string
	reader   *bufio.Reader
	writer   *bufio.Writer
}

func (client *Client) Read() {
	for {
		line, _ := client.reader.ReadString('\n')
		client.incoming <- line
	}
}

func (client *Client) Write() {
	for data := range client.outgoing {
		client.writer.WriteString(data)
		client.writer.Flush()
	}
}

func (client *Client) Listen() {
	go client.Read()
	go client.Write()
}

func NewClient(connection net.Conn) *Client {
	writer := bufio.NewWriter(connection)
	reader := bufio.NewReader(connection)

	client := &Client{
		incoming: make(chan string),
		outgoing: make(chan string),
		reader: reader,
		writer: writer,
	}

	client.Listen()

	return client
}

type ChatRoom struct {
	clients []*Client
	joins chan net.Conn
	incoming chan string
	outgoing chan string
}

func (chatRoom *ChatRoom) Broadcast(data string) {
	for _, client := range chatRoom.clients {
		client.outgoing <- data
	}
}

func (chatRoom *ChatRoom) Join(connection net.Conn) {
	client := NewClient(connection)
	chatRoom.clients = append(chatRoom.clients, client)
	go func() { for { chatRoom.incoming <- <-client.incoming } }()
}

func (chatRoom *ChatRoom) Listen() {
	go func() {
		for {
			select {
			case data := <-chatRoom.incoming:
				chatRoom.Broadcast(data)
			case conn := <-chatRoom.joins:
				chatRoom.Join(conn)
			}
		}
	}()
}

func NewChatRoom() *ChatRoom {
	chatRoom := &ChatRoom{
		clients: make([]*Client, 0),
		joins: make(chan net.Conn),
		incoming: make(chan string),
		outgoing: make(chan string),
	}

	chatRoom.Listen()

	return chatRoom
}

func main() {
	chatRoom := NewChatRoom()

	listener, _ := net.Listen("tcp", ":6666")

	for {
		conn, _ := listener.Accept()
		chatRoom.joins <- conn
	}
}

// // ----------------------------------------

/*
Simple TCP chat implementation in Go.

To start, use : go run ./server.go
By default it runs at port 6000

Connect to it via telnet -
(from own machine) : telnet localhost 6000

(from machine on same network) : telnet [local IP] 6000

(*from external machine) : telnet [Your public IP] 6000
(You need to have port forwarding on using your router)

*/
package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"io"
)

func main() {

	const maxUsers = 2 // By default

	users := make(map[net.Conn]string) // Map of active connections
	newConnection := make(chan net.Conn) // Handle new connection
	addedUser := make(chan net.Conn)   // Add new connection
	deadUser := make(chan net.Conn)    // Users that have left chat
	messages := make(chan string)      // channel that recieves messages from all users

	server, err := net.Listen("tcp", ":6000")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	go func() { // Launch routine that will accept connections
		for {
			conn, err := server.Accept()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			if len(users) < maxUsers {
				newConnection <- conn // Send to handle new user
			}else{
				io.WriteString(conn, "Server is full!")
			}
		}
	}()

	for { // Run forever

		select {
		case conn := <-newConnection:

			go func(conn net.Conn) { // Ask user for name and information
				reader := bufio.NewReader(conn)
				io.WriteString(conn, "Enter name: ")
				userName, _ := reader.ReadString('\n')
				userName = strings.Trim(userName, "\r\n")
				log.Printf("Accepted new user : %s", userName)
				messages <- fmt.Sprintf("Accepted user : [%s]\n\n", userName)

				users[conn] = userName // Add connection

				addedUser <- conn // Add user to pool
			}(conn)

		case conn := <-addedUser: // Launch a new go routine for the newly added user

			go func(conn net.Conn, userName string) {
				reader := bufio.NewReader(conn)
				for { // Run forever and handle this user's messages
					newMessage, err := reader.ReadString('\n')
					newMessage = strings.Trim(newMessage, "\r\n")
					if err != nil {
						break
					}
					// Send to messages channel therefore ring every user
					messages <- fmt.Sprintf(">%s: %s \a\n\n", userName, newMessage)
				}

				deadUser <- conn // If error occurs, connection has been terminated
				messages <- fmt.Sprintf("%s disconnected\n\n", userName)
			}(conn, users[conn])

		case message := <-messages: // If message recieved from any user

			for conn, _ := range users { // Send to all users
				go func(conn net.Conn, message string) { // Write to all user connections
						_, err := io.WriteString(conn, message)
						if err != nil {
							deadUser <- conn
						}
				}(conn, message)
				log.Printf("New message: %s", message)
				log.Printf("Sent to %d users", len(users))
			}

		case conn := <-deadUser: // Handle dead users
			log.Printf("Client disconnected")
			delete(users, conn)
		}
	}
}
