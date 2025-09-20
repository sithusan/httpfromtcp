package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"strings"
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

		lines := getLinesChannel(conn)

		for line := range lines {
			fmt.Println(line)
		}

		fmt.Println("Connection is closed from", conn.RemoteAddr())
	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {

	lines := make(chan string)
	buffer := make([]byte, 8)

	go func() {
		defer f.Close()
		defer close(lines)

		currentLine := ""

		for {
			n, err := f.Read(buffer)

			parts := strings.Split(string(buffer[:n]), "\n")

			for i := 0; i < len(parts)-1; i++ {
				currentLine = currentLine + parts[i]
				lines <- currentLine
				currentLine = ""
			}

			currentLine += parts[len(parts)-1]

			if err != nil {
				if err == io.EOF {
					break
				}

				log.Fatalf("Cannot read the file %s\n", err)
			}
		}

		if currentLine != "" {
			lines <- currentLine
		}
	}()

	return lines
}
