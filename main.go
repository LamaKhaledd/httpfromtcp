package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

const inputFilePath = "messages.txt"

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
				log.Printf("error reading file: %s\n", err)
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
	f, err := os.Open(inputFilePath)
	if err != nil {
		log.Fatalf("could not open %s: %s\n", inputFilePath, err)
	}

	for line := range getLinesChannel(f) {
		fmt.Printf("read: %s\n", line)
	}
}
