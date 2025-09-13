package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

func getLinesChannel(f io.ReadCloser) <-chan string {
	ch := make(chan string)

	go func() {
		defer f.Close()
		var currentLine string
		buffer := make([]byte, 8)

		for {
			n, err := f.Read(buffer)
			if err != nil {
				if errors.Is(err, io.EOF) {
					break
				}
				log.Printf("error reading: %s\n", err)
				break
			}

			chunk := string(buffer[:n])
			parts := strings.Split(chunk, "\n")

			for i, part := range parts {
				if i < len(parts)-1 {
					ch <- currentLine + part
					currentLine = ""
				} else {
					currentLine += part
				}
			}
		}

		if len(currentLine) > 0 {
			ch <- currentLine
		}

		close(ch)
	}()

	return ch
}

func main() {
	listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatalf("could not start listener: %s\n", err)
	}
	defer listener.Close()
	fmt.Println("TCP server listening on :42069")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		fmt.Println("Connection accepted")

		for line := range getLinesChannel(conn) {
			fmt.Println(line)
		}
		fmt.Println("Connection closed")
	}
}
