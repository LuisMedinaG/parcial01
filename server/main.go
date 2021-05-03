package main

import (
	"fmt"
	"log"
	"net"
	"encoding/gob"
	// "sync"
)

const (
	CONN_HOST = "0.0.0.0"
	CONN_PORT = "3333"
	CONN_TYPE = "tcp"
)

var connections = make(map[int]net.Conn) 
var id = 0

func main() {
	l, err := net.Listen(CONN_TYPE, CONN_HOST + ":" + CONN_PORT)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		return
	}
	
	defer l.Close()
	
	fmt.Println("Listening on " + CONN_HOST + ":" + CONN_PORT)
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
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
		// msg, err := common.ReadMsg(conn)
		if err != nil {
			log.Println(err)
			return
		}
		fmt.Printf("Message Received: %s\n", msg)
		// broadcast(conn, msg)
	}
}


func broadcast(conn net.Conn, msg string) {
	for _, c := range connections {
		if c != conn {
			err := gob.NewDecoder(c).Decode(&msg)
			if err != nil {
				log.Println(err)
			}
		}
	}
}
