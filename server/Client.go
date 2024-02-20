package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

type Client struct {
	ID     int
	Name   string
	Conn   net.Conn
	Server *Pool
	Out    chan []byte
}

func NewClient(id int, conn net.Conn) *Client {
	ret := new(Client)
	ret.ID = id
	ret.Conn = conn
	ret.Server = pool
	ret.Out = make(chan []byte, queueSize)
	return ret
}

func (c *Client) Read() {
	br := bufio.NewReader(c.Conn)
	for {
		br.Reset(c.Conn)
		buf, err := br.ReadBytes('\n')
		if err != nil {
			log.Println(err.Error())
			c.Conn.Close()
			delete(c.Server.Clients, c)
			connectedClients--
			log.Printf("Connected Clients: %d\n", connectedClients)
			out := fmt.Sprintf("*** %s left the chat ***\n", c.Name)
			log.Print(out)
			c.Server.In <- []byte(out)
			return
		}
		if c.Name == "" {
			c.Name = string(buf[:len(buf)-3])
			out := fmt.Sprintf("*** %s entered the chat ***\n", c.Name)
			log.Print(out)
			c.Server.In <- []byte(out)
		} else {
			log.Printf("%s", buf)
			c.Server.In <- buf
		}
	}
}

func (c *Client) Write() {
	for msg := range c.Out {
		_, err := c.Conn.Write(msg)
		if err != nil {
			c.Conn.Close()
			delete(c.Server.Clients, c)
			connectedClients--
			log.Printf("Connected Clients: %d\n", connectedClients)
			return
		}
	}
}
