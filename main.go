package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	file, err := os.Open("messages.txt")	
	if err != nil {
		log.Fatal(err)
	}
	
	lines := getLinesChannel(file)

	for line := range lines {
		fmt.Printf("read: %s\n", line)
	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {
	linesCh := make(chan string)

	go func(){
		defer close(linesCh)

		currentLine := "" 
		for {
			buffer := make([]byte, 8)
			_, err := f.Read(buffer)
			if err == io.EOF {
				return
			}

			parts := bytes.Split(buffer, []byte("\n"))

			for i := range(len(parts)) {
				if i != 0 {
					linesCh <- currentLine
					currentLine = ""
				}

				currentLine += string(parts[i])
			}
		}
	}()

	return linesCh
}
