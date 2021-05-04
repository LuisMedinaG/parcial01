package main

import (
	"bufio"
	"fmt"
	"encoding/gob"
	"log"
	"net"
	"io"
	"os"
)

const (
	CONN_HOST = "localhost"
	CONN_PORT = "3333"
	CONN_TYPE = "tcp"
)

func main() {
	conn, err := net.Dial(CONN_TYPE,CONN_HOST+":"+CONN_PORT)
	if err != nil {
		log.Fatal(err)
	}
	
	defer conn.Close()
	go printOutput(conn)
	sendMessage(conn)
}

func sendMessage(conn net.Conn) {
	for {
		fmt.Println("Enter text: ")
		text, err := reader.ReadString('\n')
		text = text[:len(text)-1]
		if err != nil {
			log.Fatal(err)
		}
		err = gob.NewEncoder(conn).Encode(username+" says "+text)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func printOutput(conn net.Conn) {
	for {
		var msg string
		err := gob.NewDecoder(conn).Decode(&msg)

		if err == io.EOF {
			conn.Close()
			fmt.Println("Connection Closed.")
			return
		}
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(msg)
	}
}
