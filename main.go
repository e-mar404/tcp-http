package main

import (
	"bytes"
	"fmt"
	"io"
	"net"
)

func main() {
	port := ":42069"
	listener, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Printf("error opening port: %v", err)
	}

	fmt.Printf("accepting connections on port: %s\n", port)

	for {
		conn, err := listener.Accept()

		if err != nil {
			fmt.Printf("error accepting connection: %v\n", err)
		}
		
		fmt.Printf("accepted a connection\n")

		lines := getLinesChannel(conn)
		for line := range lines {
			fmt.Printf("%s", line)
		}

		fmt.Printf("\n")
		fmt.Printf("closing connection...\n")
	}
}

func getLinesChannel(c io.ReadCloser) <-chan string {
	linesCh := make(chan string)

	go func(){
		defer close(linesCh)

		currentLine := "" 
		for {
			buffer := make([]byte, 8)
			_, err := c.Read(buffer)
			if err == io.EOF {
				break
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

		linesCh <- currentLine
	}()

	return linesCh
}
