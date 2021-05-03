package main

import (
	"bufio"
	"fmt"
	"gob"
	"log"
	"net"
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
	sendMessage(conn)
}

func sendMessage(conn *net.TCPConn) {

	fmt.Print("Enter username: ")
	reader := bufio.NewReader(os.Stdin)
	username, err := reader.ReadString('\n')
	username = username[:len(username)-1]
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Enter text: ")
	for {
		text, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		err = gob.NewEncoder(conn).Encode(msg)
		if err != nil {
			log.Fatal(err)
		}
	}
}

// func printOutput(conn *net.TCPConn) {
// 	for {

// 		msg, err := common.ReadMsg(conn)
// 		// Receiving EOF means that the connection has been closed
// 		if err == io.EOF {
// 			// Close conn and exit
// 			conn.Close()
// 			fmt.Println("Connection Closed. Bye bye.")
// 			os.Exit(0)
// 		}
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		fmt.Println(msg)
// 	}
// }
