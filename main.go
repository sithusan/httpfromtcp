package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func main() {
	file, err := os.Open("messages.txt")

	if err != nil {
		log.Fatalf("Cannot open messages %s\n", err)
	}

	defer file.Close()

	buffer := make([]byte, 8)

	currentLine := ""

	for {
		n, err := file.Read(buffer)

		parts := strings.Split(string(buffer[:n]), "\n")

		for i := 0; i < len(parts)-1; i++ {
			currentLine = currentLine + parts[i]
			fmt.Printf("read: %s\n", currentLine)
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
		fmt.Printf("read: %s\n", currentLine)
	}
}
