package main

import (
	"fmt"
	"time"
)

type Pool struct {
	In      chan []byte
	Clients map[int]*Client
}

func NewPool() *Pool {
	ret := new(Pool)
	ret.In = make(chan []byte)
	ret.Clients = make(map[int]*Client)
	return ret
}

func (p *Pool) Run() {
	for m := range p.In {
		t := time.Now().UTC().Format(time.DateTime)
		msg := fmt.Sprintf("[%s] - %s", t, m)
		for _, c := range p.Clients {
			c.Out <- []byte(msg)
		}
	}
}
