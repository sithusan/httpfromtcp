package main

import (
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	file, err := os.Open("messages.txt")

	if err != nil {
		log.Fatalf("Cannot open messages %s\n", err)
	}

	defer file.Close()

	buffer := make([]byte, 8)

	for {
		n, err := file.Read(buffer)

		if err != nil {
			if err == io.EOF {
				break
			}

			log.Fatalf("Cannot read the file %s\n", err)
		}

		fmt.Printf("read: %s\n", buffer[:n])
	}
}
