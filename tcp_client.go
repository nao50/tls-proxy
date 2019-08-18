package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
)

var (
	remoteAddress *string = flag.String("remoteAddress", "", "Remote server address")
)

func main() {
	flag.Parse()
	conn, err := net.Dial("tcp", *remoteAddress)
	if err != nil {
		fmt.Println("ERROR", err)
		os.Exit(1)
	}

	userInput := bufio.NewReader(os.Stdin)
	response := bufio.NewReader(conn)
	for {
		fmt.Print("Client Input  > ")
		userLine, err := userInput.ReadBytes(byte('\n'))
		switch err {
		case nil:
			conn.Write(userLine)
		case io.EOF:
			os.Exit(0)
		default:
			fmt.Println("ERROR", err)
			os.Exit(1)
		}

		serverLine, err := response.ReadBytes(byte('\n'))
		switch err {
		case nil:
			fmt.Printf("Server Output > %s", string(serverLine))
		case io.EOF:
			os.Exit(0)
		default:
			fmt.Println("ERROR", err)
			os.Exit(2)
		}
	}
}
