package main

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
		for _, c := range p.Clients {
			c.Out <- m
		}
	}
}
