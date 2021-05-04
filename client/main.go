package main

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
)

const (
	HOST = "localhost"
	PORT = "3333"
)

type Client struct {
	userName string
	msg      string
	file     []byte
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

var messages = make([]string, 0)

func leerMensajesLlegada(conn net.Conn) {
	for {
		var msg string
		err := gob.NewDecoder(conn).Decode(&msg)

		if err == io.EOF {
			conn.Close()
			fmt.Println("Serve closed connection.")
			os.Exit(0)
		}
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}

		messages = append(messages, msg)
		fmt.Println(msg)
	}
}

func main() {
	conn, err := net.Dial("tcp", HOST+":"+PORT)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer conn.Close()

	var userName string
	firstMsg := true

	go leerMensajesLlegada(conn)
	for {
		opt := getMenuOpt()
		switch opt {
		case 1:
			if firstMsg {
				reader := bufio.NewReader(os.Stdin)
				fmt.Print("Usuario: ")
				userName, _ = reader.ReadString('\n')
				userName = userName[:len(userName)-1]
				firstMsg = false
			}
			reader := bufio.NewReader(os.Stdin)
			prefix := "@" + userName + ": "
			fmt.Print(prefix)
			text, _ := reader.ReadString('\n')
			text = text[:len(text)-1]
			text = prefix + text

			err := gob.NewEncoder(conn).Encode(text)
			if err != nil {
				log.Fatal(err)
			}
		case 2:

		case 3:
			for msg := range messages {
				fmt.Println("@" + userName + ": " + strconv.Itoa(msg))
			}
		case 4:
			log.Println("Client terminated.")
			return
		}

	}
}
