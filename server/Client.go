package main

import (
	"bufio"
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

func (c *Client) Write() {
	for msg := range c.Out {
		n, err := c.Conn.Write(msg)
		if err != nil {
			log.Fatalf("%s: bytes written: %d", err.Error(), n)
		}
	}
}

func (c *Client) Read() {
	br := bufio.NewReader(c.Conn)
	for {
		br.Reset(c.Conn)
		buf, err := br.ReadBytes('\n')
		if err != nil {
			log.Fatal(err.Error())
		}
		c.Server.In <- buf
	}
}
