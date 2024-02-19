package main

import (
	"fmt"
	"log"
	"net"
)

const (
	network = "tcp4"
	addr    = "127.0.0.1:49000"
)

func main() {
	listener, err := net.Listen(network, addr)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer listener.Close()
	fmt.Println("Server listening for connections...")
	conn, err := listener.Accept()
	if err != nil {
		log.Fatal(err.Error())
	}
	defer conn.Close()
	fmt.Println("Client connected.")
}
