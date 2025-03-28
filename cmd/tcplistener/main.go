package main

import (
	"errors"
	"fmt"
	"io"
	"net"
	"strings"
)

func main() {
	port := ":42069"
	listener, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Printf("error opening port: %v", err)
	}

	fmt.Printf("accepting connections on localhost%s\n", port)
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("error accepting connection: %v\n", err)
		}
		fmt.Printf("accepted a connection\n")

		lines := getLinesChannel(conn)
		for line := range lines {
			fmt.Printf("%s\n", line)
		}

		fmt.Printf("closing %v\n", conn.RemoteAddr())
	}
}

func getLinesChannel(c io.ReadCloser) <-chan string {
	linesCh := make(chan string)

	go func() {
		defer c.Close()
		defer close(linesCh)

		currentLine := ""
		for {
			buffer := make([]byte, 8)
			n, err := c.Read(buffer)
			if err != nil {
				if currentLine != "" {
					linesCh <- currentLine
				}
				if errors.Is(err, io.EOF) {
					break
				}
			}

			str := string(buffer[:n])
			parts := strings.Split(str, "\n")
			for i := range len(parts) - 1 {
				linesCh <- fmt.Sprintf("%s%s", currentLine, parts[i])
				currentLine = ""
			}

			currentLine += parts[len(parts)-1]
		}
	}()

	return linesCh
}
