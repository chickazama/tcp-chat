package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

type Client struct {
	ID     int
	Conn   net.Conn
	Server *Pool
	Out    chan []byte
}

func NewClient(id int, conn net.Conn) *Client {
	ret := new(Client)
	ret.ID = id
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
			c.Server.Logout <- c.ID
			return
		}
		fmt.Printf("%s", buf)
		c.Server.In <- buf
	}
}

func (c *Client) Write() {
	for msg := range c.Out {
		_, err := c.Conn.Write(msg)
		if err != nil {
			log.Println("Client exited")
			c.Server.Logout <- c.ID
			return
		}
	}
}
