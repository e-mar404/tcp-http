package request

import (
	"fmt"
	"io"
	"strings"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	Method        string
	RequestTarget string
	HttpVersion   string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	buffer, err := io.ReadAll(reader)
	if err != nil {
		fmt.Printf("error reading from reader: %s\n", err)
		return nil, err
	}

	str := string(buffer)
	parts := strings.Split(str, "\r\n")

	requestLineStr := strings.Split(parts[0], " ")
	
	return &Request{
		RequestLine: RequestLine{
			Method: requestLineStr[0],
			RequestTarget: requestLineStr[1],
			HttpVersion: requestLineStr[2],
		},
	}, nil

}
