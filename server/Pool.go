package main

import (
	"fmt"
	"log"
	"time"
)

const (
	queueSize = 8
)

type Pool struct {
	In      chan []byte
	Clients map[*Client]bool
}

func NewPool() *Pool {
	ret := new(Pool)
	ret.In = make(chan []byte, queueSize)
	ret.Clients = make(map[*Client]bool)
	return ret
}

func (p *Pool) Run() {
	for m := range p.In {
		t := time.Now().UTC().Format(time.DateTime)
		msg := fmt.Sprintf("[%s] - %s", t, m)
		out := []byte(msg)
		mtx.Lock()
		_, err := fp.Write([]byte(out))
		if err != nil {
			log.Println(err.Error())
		}
		mtx.Unlock()
		out[len(out)-1] = 0
		for c := range p.Clients {
			c.Out <- []byte(out)
		}
	}
}
