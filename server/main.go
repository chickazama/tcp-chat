package main

import (
	"bufio"
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
	br := bufio.NewReader(conn)
	for {
		br.Reset(conn)
		inBuf, err := br.ReadBytes('\n')
		if err != nil {
			log.Fatal(err.Error())
		}
		fmt.Printf("%s", string(inBuf))
		n, err := conn.Write([]byte("Server: OK!\n"))
		if err != nil {
			log.Fatalf("%s: bytes written: %d", err.Error(), n)
		}
	}
}
