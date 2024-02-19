package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

const (
	network = "tcp4"
	addr    = "127.0.0.1:49000"
)

func main() {
	conn, err := net.Dial(network, addr)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer conn.Close()
	go read(conn)
	write(conn)
}

func read(conn net.Conn) {
	br := bufio.NewReader(os.Stdin)
	for {
		br.Reset(os.Stdin)
		outBuf, err := br.ReadBytes('\n')
		if err != nil {
			log.Fatal(err.Error())
		}
		n, err := conn.Write(outBuf)
		if err != nil {
			log.Fatalf("%s: bytes written: %d", err.Error(), n)
		}
	}
}

func write(conn net.Conn) {
	cr := bufio.NewReader(conn)
	for {
		cr.Reset(conn)
		inBuf, err := cr.ReadBytes('\n')
		if err != nil {
			log.Fatal(err.Error())
		}
		fmt.Printf("%s\n", string(inBuf))
	}
}
