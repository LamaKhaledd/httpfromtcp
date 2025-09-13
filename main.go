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

func main() {
	f, err := os.Open(inputFilePath)
	if err != nil {
		log.Fatalf("could not open %s: %s\n", inputFilePath, err)
	}
	defer f.Close()

	var currentLine string

	for {
		b := make([]byte, 8)
		n, err := f.Read(b)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			fmt.Printf("error: %s\n", err.Error())
			return
		}

		chunk := string(b[:n])
		parts := strings.Split(chunk, "\n")

		for i, part := range parts {
			if i < len(parts)-1 {
				fmt.Printf("read: %s\n", currentLine+part)
				currentLine = ""
			} else {
				currentLine += part
			}
		}
	}

	if len(currentLine) > 0 {
		fmt.Printf("read: %s\n", currentLine)
	}
}
