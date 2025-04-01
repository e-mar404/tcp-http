package main

import (
	"e-mar404/httpfromtcp/internal/request"
	"fmt"
	"net"
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

		req, err := request.RequestFromReader(conn)
		if err != nil {
			fmt.Printf("error reading request: %v\n", err)
		}
		
		printRequest(req)

		fmt.Printf("closing %v\n", conn.RemoteAddr())
	}
}

func printRequest(r *request.Request) {
	fmt.Println("Request line:")
	fmt.Printf("- Method: %s\n", r.RequestLine.Method)
	fmt.Printf("- Target: %s\n", r.RequestLine.RequestTarget)
	fmt.Printf("- Version: %s\n", r.RequestLine.HttpVersion)
}
