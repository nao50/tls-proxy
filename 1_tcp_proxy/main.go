package main

import (
	"fmt"
	"io"
	"net"
	"sync"
)

type Proxy struct {
	from string
	to   string
	done chan struct{}
}

func NewProxy(from, to string) *Proxy {
	return &Proxy{
		from: from,
		to:   to,
		done: make(chan struct{}),
	}
}

func (p *Proxy) Start() error {
	listener, err := net.Listen("tcp", p.from)
	if err != nil {
		return err
	}
	go p.run(listener)
	return nil
}

func (p *Proxy) Stop() {
	if p.done == nil {
		return
	}
	close(p.done)
	p.done = nil
}

func (p *Proxy) run(listener net.Listener) {
	for {
		select {
		case <-p.done:
			return
		default:
			connection, err := listener.Accept()
			if err == nil {
				go p.handle(connection)
			}
		}
	}
}

func (p *Proxy) handle(connection net.Conn) {
	defer connection.Close()
	remote, err := net.Dial("tcp", p.to)
	if err != nil {
		return
	}
	defer remote.Close()
	wg := &sync.WaitGroup{}
	wg.Add(2)
	go p.copy(remote, connection, wg)
	go p.copy(connection, remote, wg)
	wg.Wait()
}

func (p *Proxy) copy(from, to net.Conn, wg *sync.WaitGroup) {
	defer wg.Done()
	select {
	case <-p.done:
		return
	default:
		if _, err := io.Copy(to, from); err != nil {
			p.Stop()
			return
		}
	}
}

func main() {
	proxy := NewProxy("localhost:8080", "localhost:8888")

	var wg sync.WaitGroup
	wg.Add(1)

	err := proxy.Start()
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	wg.Wait()

}
