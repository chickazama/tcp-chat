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

var (
	msg chan []byte
)

func init() {
	msg = make(chan []byte)
}

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
	go read(conn)
	write(conn)
}

func read(conn net.Conn) {
	defer conn.Close()
	br := bufio.NewReader(conn)
	for {
		br.Reset(conn)
		buf, err := br.ReadBytes('\n')
		if err != nil {
			log.Fatal(err.Error())
		}
		fmt.Printf("%s", string(buf))
		msg <- buf
	}
}

func write(conn net.Conn) {
	for m := range msg {
		n, err := conn.Write(m)
		if err != nil {
			log.Fatalf("%s: bytes written: %d", err.Error(), n)
		}
	}
}
