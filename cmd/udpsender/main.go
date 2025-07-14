package main

import (
	"bufio"
	"log"
	"net"
	"os"
	"os/signal"
)

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		os.Exit(0)
	}()

	addr, err := net.ResolveUDPAddr("udp", ":42069")
	if err != nil {
		log.Fatalf("Error resolving UDP address: %v", err)
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		log.Fatalf("Failed to establish a UDP connection: %v", err)
	}
	defer func() { _ = conn.Close() }()

	reader := bufio.NewReader(os.Stdin)

	for {
		log.Println(">")

		buf, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf("Error reading string: %v", err)
		}

		_, err = conn.Write([]byte(buf))
		if err != nil {
			log.Fatalf("Error writing bytes: %v", err)
		}
	}
}
