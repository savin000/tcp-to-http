package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
)

func main() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		os.Exit(0)
	}()

	listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = listener.Close() }()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Connection Accepted")

		linesChannel := getLinesChannel(conn)

		for line := range linesChannel {
			fmt.Printf("%s\n", line)
		}

		fmt.Println("Connection Closed")
	}
}

func getLinesChannel(conn net.Conn) <-chan string {
	result := make(chan string)
	chunkSize := 8
	chunk := make([]byte, chunkSize)
	currentLine := ""

	go func() {
		defer func() { _ = conn.Close() }()
		defer close(result)

		for {
			chunkBytes, err := conn.Read(chunk)

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
