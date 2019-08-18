package main

import (
	"bufio"
	"fmt"
	"net"
)

func echo(conn net.Conn) {
	r := bufio.NewReader(conn)
	for {
		line, err := r.ReadBytes(byte('\n'))
		if err != nil {
			fmt.Println("BREAK: ", err)
			break
		}
		fmt.Printf("%v > %s", conn.RemoteAddr(), line)
		conn.Write(line)
	}
}

func main() {
	listen, err := net.Listen("tcp", "localhost:8888")
	if err != nil {
		fmt.Println("ERROR: ", err)
		return
	}
	fmt.Println(">> Server Running localhost:8888")

	for {
		conn, err := listen.Accept()
		fmt.Printf("conn: %v \n", conn.RemoteAddr())
		if err != nil {
			fmt.Println("Accept Error: ", err)
			continue
		}
		go echo(conn)
	}
}
