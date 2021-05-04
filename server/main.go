package main

import (
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"net"
)

const (
	HOST = "localhost"
	PORT = "3333"
)

type Message struct {
	Text string
	File []byte
}

type Chat struct {
	usrCount int
	Host     string
	Port     string
	clients  map[net.Conn]int
	messages []string
}

func NewChat(host string, port string) *Chat {
	chat := &Chat{
		usrCount: 0,
		Host:     host,
		Port:     port,
		clients:  make(map[net.Conn]int),
		messages: make([]string, 0),
	}
	go chat.Listen()

	return chat
}

func (chat *Chat) Listen() {
	listener, err := net.Listen("tcp", chat.Host+":"+chat.Port)
	if err != nil {
		log.Println(err)
		return
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			return
		}

		chat.clients[conn] = chat.usrCount
		chat.usrCount++

		go chat.HandleRequest(conn)
	}
}

func (chat *Chat) HandleRequest(conn net.Conn) {
	defer func() {
		conn.Close()
		delete(chat.clients, conn)
	}()

	var text string
	var message Message
	for {
		err := gob.NewDecoder(conn).Decode(&message)
		if err == io.EOF {
			// log.Println("Client connection closed.")
			return
		}
		if err != nil {
			log.Fatal("ERROR Listening: ", err)
			return
		}
		if len(message.File) > 0 {
			// TODO recive file
			log.Println("File recived")
			// save file
			// messages = append(messages, message.File)
			// text = message.FileName
		} else {
			text = message.Text
		}
		chat.messages = append(chat.messages, text)
		chat.Broadcast(message)
	}
}

func (chat *Chat) Broadcast(message Message) {
	for conn := range chat.clients {
		err := gob.NewEncoder(conn).Encode(&message)
		if err != nil {
			log.Fatal("Error Broadcasting:", err)
		}
	}
}

func getMenuOptServ() (opt int) {
	fmt.Println(`
********* CLIENTE *********
1. Mostrar mensajes/archivos enviados
2. Respaldar mensajes/archivos enviados
3. Terminar servidor
Ingrese opcion: `)
	fmt.Scan(&opt)
	return opt
}

func main() {
	chat := NewChat(HOST, PORT)

	for {
		opt := getMenuOptServ()
		switch opt {
		case 1:
			fmt.Println("\n ---- Mostrar mensajes ----\n ")
			for _, msg := range chat.messages {
				fmt.Println(msg)
			}
		case 2:
			fmt.Println("\n ---- Respaldar mensaje ----\n ")
			// TODO respaldar lista de mensaje en un archivo
		case 3:
			fmt.Println("\nServidor terminado.")
			return
		default:
			fmt.Println("\nERROR: Opcion invalida.")
		}
	}
}
