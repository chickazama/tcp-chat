package core

import (
	"bufio"
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

func New(conn net.Conn) *Client {
	ret := new(Client)
	ret.Conn = conn
	ret.Send = make(chan []byte, queueSize)
	ret.Receive = make(chan []byte)
	return ret
}

func (c *Client) Read() {
	br := bufio.NewReader(c.Conn)
	for {
		buf, err := br.ReadBytes(0)
		if err != nil {
			// log.Println(err.Error())
			return
		}
		buf[len(buf)-1] = '\n'
		c.Receive <- buf
		br.Reset(c.Conn)
	}
}
