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
	HOST = "localhost"
	PORT = "3333"
)

// type File struct {
// 	Bytes    []byte
// 	fileName string
// 	ext      string
// }

type Client struct {
	conn     net.Conn
	host     string
	port     string
	usrName  string
	messages []string
}

type Message struct {
	Text string
	File []byte
}

func NewClient(host string, port string) *Client {
	usrName := Input("Usuario: ")
	client := &Client{
		usrName:  usrName,
		host:     host,
		port:     port,
		messages: make([]string, 0),
	}

	client.Dial()

	return client
}

func (client *Client) Dial() {
	conn, err := net.Dial("tcp", client.host+":"+client.port)
	if err != nil {
		log.Fatal(err)
		return
	}

	client.conn = conn
	go client.HandleServerMsg()
}

func (client *Client) HandleServerMsg() {
	defer client.conn.Close()

	var text string
	var message Message
	for {
		err := gob.NewDecoder(client.conn).Decode(&message)
		if err == io.EOF {
			fmt.Println("Server connection closed.")
			os.Exit(0)
		}
		if err != nil {
			log.Fatal("ERROR Listening: ", err)
			os.Exit(1)
		}
		if len(message.File) > 0 {
			log.Println("File recived")
			// messages = append(messages, message.File)
			// text = message.FileName
		} else {
			text = message.Text
		}
		client.messages = append(client.messages, text)
	}
}

func (client *Client) sendMessge(msg Message) {
	err := gob.NewEncoder(client.conn).Encode(msg)
	if err != nil {
		log.Fatal(err)
	}
}

func Input(inputTxt string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(inputTxt)
	text, _ := reader.ReadString('\n')
	return text[:len(text)-1]
}

func getMenuOpt() (opt int) {
	fmt.Println(`
********* CLIENTE *********

1. Enviar mensaje
2. Enviar archivo
3. Mostrar mensajes recibidos
4. Terminar cliente
Ingrese opcion: `)
	fmt.Scan(&opt)
	return opt
}

func main() {
	client := NewClient(HOST, PORT)

	var message Message
	for {
		opt := getMenuOpt()
		switch opt {
		case 1:
			fmt.Println("\n ---- Enviar mensaje ----\n")

			text := Input(">> ")
			message := Message{
				Text: "@" + client.usrName + ": " + text,
			}
			client.sendMessge(message)
		case 2:
			fmt.Println("\n --- Enviar archivo ---\n")
			message.File = make([]byte, 5)

			client.sendMessge(message)
		case 3:
			fmt.Println("\n ---- Mensajes chat ----\n")
			for _, msg := range client.messages {
				fmt.Println(msg)
			}
		case 4:
			fmt.Println("\nConexion terminada.")
			return
		default:
			fmt.Println("\nERROR: Opcion invalida.")

		}
	}
}
