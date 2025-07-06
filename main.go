package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func main() {
	file, err := os.Open("message.txt")
	if err != nil {
		log.Fatal(err)
	}

	linesChannel := getLinesChannel(file)

	for line := range linesChannel {
		fmt.Printf("read: %s\n", line)
	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {
	result := make(chan string)
	chunkSize := 8
	chunk := make([]byte, chunkSize)
	currentLine := ""

	go func() {
		defer func() { _ = f.Close() }()
		defer close(result)

		for {
			chunkBytes, err := f.Read(chunk)

			if err != nil {
				if err == io.EOF {
					if len(currentLine) > 0 {
						result <- currentLine
					}
					break
				}
				log.Fatalf("Error reading file: %v", err)
			}

			parts := strings.Split(string(chunk[:chunkBytes]), "\n")

			for i, part := range parts {
				if i < len(parts)-1 {
					currentLine += part
					result <- currentLine
					currentLine = ""
				}
			}
			currentLine += parts[len(parts)-1]
		}
	}()

	return result
}
