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

var (
	pool *Pool
)

func init() {
	pool = NewPool()
	go pool.Run()
}

func main() {
	listener, err := net.Listen(network, addr)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer listener.Close()
	fmt.Println("Server listening for connections...")
	for i := 0; i < 10; i++ {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err.Error())
		}
		c := NewClient(conn)
		pool.Clients[i] = c
		go c.Read()
		go c.Write()
	}
}
