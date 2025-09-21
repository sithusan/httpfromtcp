package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

const host = "localhost"
const port = "42069"

func main() {
	udpAddr, err := net.ResolveUDPAddr("udp", net.JoinHostPort(host, port))

	if err != nil {
		log.Fatalf("Cannot resolve udp traffic %s\n", err)
	}

	conn, err := net.DialUDP("udp", nil, udpAddr)

	if err != nil {
		log.Fatalf("Cannot dail UDP %s", err)
	}

	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println(">")

		data, err := reader.ReadString('\n')

		if err != nil {
			log.Fatalf("Cannot read %s", err)
		}

		conn.Write([]byte(data))
	}
}
