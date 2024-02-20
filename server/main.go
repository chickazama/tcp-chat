package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"sync"
)

const (
	network = "tcp4"
	addr    = "127.0.0.1:49000"
)

const (
	maxConnections = 10
)

var (
	mtx              sync.Mutex
	nextId           = 1
	connectedClients int
	pool             *Pool
	fp               *os.File
)

func init() {
	var err error
	fp, err = os.OpenFile("history.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err.Error())
	}
	pool = NewPool()
	go pool.Run()
}

func main() {
	defer fp.Close()
	listener, err := net.Listen(network, addr)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer listener.Close()
	fmt.Println("Server listening for connections...")
	for connectedClients < maxConnections {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err.Error())
		}
		c := NewClient(nextId, conn)
		pool.Clients[c] = true
		connectedClients++
		log.Printf("Connected clients: %d\n", connectedClients)
		mtx.Lock()
		buf, err := os.ReadFile("history.txt")
		if err != nil {
			log.Println(err.Error())
			buf = []byte("Chat History Unavailable\n")
		}
		if len(buf) <= 0 {
			str := "*** TCP CHAT ***\n"
			buf = []byte(str)
			fp.Write(buf)
		}
		mtx.Unlock()
		buf[len(buf)-1] = 0
		c.Out <- buf
		go c.Read()
		go c.Write()
	}
}
