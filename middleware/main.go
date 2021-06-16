package main
import (
	// "encoding/gob"
	// "errors"
	"fmt"
	"net"
	// "net/http"
	"net/rpc"
	"log"
)

type Chat struct {
	Topic	string
	usrCount int
	Host     string
	Port     string
	clients  map[net.Conn]int
	messages []string
}

type Midleware struct {
	Host     string
	Port     string
}

func (midleware *Midleware) GetAvailablePorts(args bool, reply *string) error  {
	for port, taken := range AvailablePorts {
		if taken == false {
			AvailablePorts[port] = true
			*reply = port
			return nil
		}
	}
	return fmt.Errorf("ERROR: no hay puertos disponibles")
}

var Chats map[string]Chat
var AvailablePorts map[string]bool


func main()  {
	AvailablePorts = make(map[string]bool)
	
	api := new(Midleware)
	err := rpc.Register(api)
	if err != nil {
		log.Fatal("Error registrando API", err)
	}

	fmt.Println("[INFO] API registrada")
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error arrancando servidor", err)
	}

	fmt.Println("[INFO] Arrancando servidor...")
	for {
		c, err := ln.Accept()
		if err != nil {
			fmt.Println("Error aceptando conexion", err)
			continue
		}

		go rpc.ServeConn(c)
	}
}
