package main

import (
	"fmt"
	"time"
)

type Pool struct {
	In      chan []byte
	Logout  chan int
	Clients map[int]*Client
}

func NewPool() *Pool {
	ret := new(Pool)
	ret.In = make(chan []byte)
	ret.Logout = make(chan int)
	ret.Clients = make(map[int]*Client)
	return ret
}

func (p *Pool) Run() {
	for {
		select {
		case m := <-p.In:
			t := time.Now().UTC().Format(time.DateTime)
			msg := fmt.Sprintf("[%s] - %s", t, m)
			for _, c := range p.Clients {
				c.Out <- []byte(msg)
			}
		case id := <-p.Logout:
			c, exists := p.Clients[id]
			if exists {
				c.Conn.Close()
				delete(p.Clients, id)
			}
			for _, c := range p.Clients {
				c.Out <- []byte("Client exited.\n")
			}
		}
	}
}
