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

	lines := getLinesChannel(file)

	for line := range lines {
		fmt.Printf("read: %s\n", line)
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
