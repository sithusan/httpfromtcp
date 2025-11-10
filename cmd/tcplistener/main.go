package main

import (
	"fmt"
	"log"
	"net"

	"github.com/sithusan/httpfromtcp/internal/request"
)

const port = ":42069"

func main() {
	tcpListener, err := net.Listen("tcp", port)

	if err != nil {
		log.Fatalf("Cannot listen tcp traffic %s\n", err)
	}

	defer tcpListener.Close()

	fmt.Println("Listening TCP traffic on", port)

	for {
		conn, err := tcpListener.Accept()

		if err != nil {
			log.Fatalf("Cannot accept connection %s\n", err)
		}

		fmt.Println("Connection is accepted from", conn.RemoteAddr())

		r, err := request.RequestFromReader(conn)

		if err != nil {
			log.Fatalf("Error parsing request: %s\n", err.Error())
		}

		fmt.Println("Request line:")
		fmt.Printf("- Method: %s\n", r.RequestLine.Method)
		fmt.Printf("- Target: %s\n", r.RequestLine.RequestTarget)
		fmt.Printf("- Version: %s\n", r.RequestLine.HttpVersion)

		fmt.Println("Connection is closed from", conn.RemoteAddr())
	}
}
