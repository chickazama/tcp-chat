package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

type Client struct {
	Conn   net.Conn
	Server *Pool
	Out    chan []byte
}

func NewClient(conn net.Conn) *Client {
	ret := new(Client)
	ret.Conn = conn
	ret.Server = pool
	ret.Out = make(chan []byte)
	return ret
}

func (c *Client) Read() {
	br := bufio.NewReader(c.Conn)
	for {
		br.Reset(c.Conn)
		buf, err := br.ReadBytes('\n')
		if err != nil {
			log.Println("Client exited")
			c.Conn.Close()
		}
		fmt.Printf("%s", buf)
		c.Server.In <- buf
	}
}

func (c *Client) Write() {
	for msg := range c.Out {
		n, err := c.Conn.Write(msg)
		if err != nil {
			log.Fatalf("%s: bytes written: %d", err.Error(), n)
		}
	}
}
