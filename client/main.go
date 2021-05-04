package main

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

const (
	HOST      = "localhost"
	PORT      = "3333"
)
var messages []string

type Client struct {
	messages []string
	conn 	net.Conn
}

type Message struct {
	Text     string
	File     []byte
}


func getMenuOpt() (opt int) {
	fmt.Print(`
********* CLIENTE *********

1. Enviar mensaje
2. Enviar archivo
3. Mostrar mensajes recibidos
4. Terminar cliente
Ingrese opcion: `)
	fmt.Scan(&opt)
	return opt
}

func HandleRequest(conn net.Conn) {
	var message Message
	for {
		err := gob.NewDecoder(conn).Decode(&message)
		if err == io.EOF {
			fmt.Println("Server connection closed.")
			os.Exit(0)
		}
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
		messages = append(messages, message.Text)
	}
}

func main() {
	conn, err := net.Dial("tcp", HOST+":"+PORT)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer conn.Close()

	var message Message 
	var usrSet bool
	messages = make([]string, 0)
	
	go HandleRequest(conn)
	for {
		opt := getMenuOpt()
		switch opt {
		case 1:
			fmt.Println("\n ---- Enviar mensaje ----\n")
			if !usrSet {
				reader := bufio.NewReader(os.Stdin)
				fmt.Print("Usuario: ")
				userName, _ := reader.ReadString('\n')
				UserName = userName[:len(userName)-1]
				usrSet = true
			}
			reader := bufio.NewReader(os.Stdin)
			fmt.Print(">> ")
			msg, _ := reader.ReadString('\n')
			message.Text  = "@" + UserName + ": " +  msg[:len(msg)-1]
			
			err := gob.NewEncoder(conn).Encode(message)
			if err != nil {
				log.Fatal(err)
			}
		case 2:
			fmt.Println("\n --- Enviar archivo ---\n")
			message.File = make([]byte, 5, 5)
		case 3:
			fmt.Println("\n ---- Mensajes chat ----\n")
			for _, msg := range messages {
				fmt.Println(msg)
			}
		case 4:
			fmt.Println("\nConexion terminada.\n")
			return
		}
	}
}
