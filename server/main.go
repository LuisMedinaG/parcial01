package main

import (
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"net"
)

const (
	HOST      = "localhost"
	PORT      = "3333"
)

type Message struct {
	Text     string
	File     []byte
}

type ChatRoom struct {
	clients  [net.Conn]string
	messages []string
}

func NewChatRoom() *ChatRoom {
	chatRoom := &ChatRoom{
		connections = make(map[net.Conn]string, 0)
		messages = make([]string, 0)
	}

	chatRoom.Listen()

	return chatRoom
}

func (chatRoom *ChatRoom) Listen() {
	listener, err := net.Listen("tcp", HOST+":"+PORT)
	log.Println("Listening on " + HOST + ":" + PORT)
	if err != nil {
		log.Println(err)
		return
	}
	defer listener.Close()

	// go routine here or...
	for {
		conn, _ := listener.Accept()
		go HandleRequest(conn)
		// menu here?
	}
}

func (chatRoom *ChatRoom) HandleRequest(conn net.Conn) {
	defer func(){
		conn.Close()
		delete(chatRoom.clients, conn)
	}()
		
	chatRoom.clients[conn] = "id"
	var message Message
	for {
		err := gob.NewDecoder(conn).Decode(&message)
		if err == io.EOF {
			fmt.Println("Client connection closed.")
			return
		}		
		if err != nil {
			log.Println(err)
			return
		}
		
		if len(message.File) > 0 {
			log.Println("File recived")
			// messages = append(messages, message.File)
		} else {
			messages = append(messages, message.Text)
		}
		Broadcast(message)
	}
}

func (chatRoom *ChatRoom) Broadcast(message Message) {
	for conn, _ := range chatRoom.clients {
		err := gob.NewEncoder(conn).Encode(&message)
		if err != nil {
			log.Fatal("Error Broadcasting:", err)
		}
	}
}

func getMenuOptServ() (opt int) {
	fmt.Print(`
********* CLIENTE *********
1. Mostrar mensajes/archivos enviados
2. Respaldar mensajes/archivos enviados
3. Terminar servidor
Ingrese opcion: `)
	fmt.Scan(&opt)
	return opt
}

func main() {
	// chatRoom := NewChatRoom()
	_ := NewChatRoom()
}


// // ----------------------------------------

// /*
// Simple TCP chat implementation in Go.

// To start, use : go run ./server.go
// By default it runs at port 6000

// Connect to it via telnet -
// (from own machine) : telnet localhost 6000

// (from machine on same network) : telnet [local IP] 6000

// (*from external machine) : telnet [Your public IP] 6000
// (You need to have port forwarding on using your router)

// */
// package main

// import (
// 	"bufio"
// 	"fmt"
// 	"log"
// 	"net"
// 	"os"
// 	"strings"
// 	"io"
// )

// func main() {

// 	const maxUsers = 2 // By default

// 	users := make(map[net.Conn]string) // Map of active connections
// 	newConnection := make(chan net.Conn) // Handle new connection
// 	addedUser := make(chan net.Conn)   // Add new connection
// 	deadUser := make(chan net.Conn)    // Users that have left chat
// 	messages := make(chan string)      // channel that recieves messages from all users

// 	server, err := net.Listen("tcp", ":6000")
// 	if err != nil {
// 		fmt.Println(err)
// 		os.Exit(1)
// 	}

// 	go func() { // Launch routine that will accept connections
// 		for {
// 			conn, err := server.Accept()
// 			if err != nil {
// 				fmt.Println(err)
// 				os.Exit(1)
// 			}
// 			if len(users) < maxUsers {
// 				newConnection <- conn // Send to handle new user
// 			}else{
// 				io.WriteString(conn, "Server is full!")
// 			}
// 		}
// 	}()

// 	for { // Run forever

// 		select {
// 		case conn := <-newConnection:

// 			go func(conn net.Conn) { // Ask user for name and information
// 				reader := bufio.NewReader(conn)
// 				io.WriteString(conn, "Enter name: ")
// 				userName, _ := reader.ReadString('\n')
// 				userName = strings.Trim(userName, "\r\n")
// 				log.Printf("Accepted new user : %s", userName)
// 				messages <- fmt.Sprintf("Accepted user : [%s]\n\n", userName)

// 				users[conn] = userName // Add connection

// 				addedUser <- conn // Add user to pool
// 			}(conn)

// 		case conn := <-addedUser: // Launch a new go routine for the newly added user

// 			go func(conn net.Conn, userName string) {
// 				reader := bufio.NewReader(conn)
// 				for { // Run forever and handle this user's messages
// 					newMessage, err := reader.ReadString('\n')
// 					newMessage = strings.Trim(newMessage, "\r\n")
// 					if err != nil {
// 						break
// 					}
// 					// Send to messages channel therefore ring every user
// 					messages <- fmt.Sprintf(">%s: %s \a\n\n", userName, newMessage)
// 				}

// 				deadUser <- conn // If error occurs, connection has been terminated
// 				messages <- fmt.Sprintf("%s disconnected\n\n", userName)
// 			}(conn, users[conn])

// 		case message := <-messages: // If message recieved from any user

// 			for conn, _ := range users { // Send to all users
// 				go func(conn net.Conn, message string) { // Write to all user connections
// 						_, err := io.WriteString(conn, message)
// 						if err != nil {
// 							deadUser <- conn
// 						}
// 				}(conn, message)
// 				log.Printf("New message: %s", message)
// 				log.Printf("Sent to %d users", len(users))
// 			}

// 		case conn := <-deadUser: // Handle dead users
// 			log.Printf("Client disconnected")
// 			delete(users, conn)
// 		}
// 	}
// }
