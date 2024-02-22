package core

import (
	"bufio"
	"log"
	"net"
)

const (
	queueSize = 8
	network   = "tcp4"
	addr      = "127.0.0.1:49000"
)

type Client struct {
	Conn    net.Conn
	Send    chan []byte
	Receive chan []byte
}

func New() *Client {
	conn, err := net.Dial(network, addr)
	if err != nil {
		log.Fatal(err.Error())
	}
	ret := new(Client)
	ret.Conn = conn
	ret.Send = make(chan []byte, queueSize)
	ret.Receive = make(chan []byte)
	go ret.read()
	go ret.write()
	return ret
}

func (c *Client) read() {
	br := bufio.NewReader(c.Conn)
	for {
		buf, err := br.ReadBytes(0)
		if err != nil {
			log.Println(err.Error())
			return
		}
		buf[len(buf)-1] = '\n'
		c.Receive <- buf
		br.Reset(c.Conn)
	}
}

func (c *Client) write() {
	for buf := range c.Send {
		n, err := c.Conn.Write(buf)
		if err != nil {
			log.Printf("Bytes Written: %d: %s\n", n, err.Error())
			return
		}
	}
}

// func (c *Client) handleMessage(buf []byte) {
// }
