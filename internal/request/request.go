package request

import (
	"fmt"
	"io"
	"regexp"
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

const crlf = "\r\n"

func RequestFromReader(reader io.Reader) (*Request, error) {
	rawBytes, err := io.ReadAll(reader)
	if err != nil {
		fmt.Printf("error reading from reader: %s\n", err)
		return nil, err
	}

	requestLine, err := parseRequestLine(rawBytes)
	if err != nil {
		return nil, err
	}

	return &Request{
		RequestLine: *requestLine,
	}, nil
}

func parseRequestLine(data []byte) (*RequestLine, error) {
	parts := strings.Split(string(data), crlf)

	requestLineParts := strings.Split(parts[0], " ")

	if len(requestLineParts) != 3 {
		return nil, fmt.Errorf("Expecting 3 parts to the Request Line, got %v", len(requestLineParts))
	}

	requestLinePatters := []*regexp.Regexp{
		regexp.MustCompile(`[A-Z]{3,}`),
		regexp.MustCompile(`/\w*`),
		regexp.MustCompile(`HTTP/1.1`),
	}

	for i, pattern := range requestLinePatters {
		if !pattern.MatchString(requestLineParts[i]) {
			return nil, fmt.Errorf("Request Line Malformed")
		}
	}
	
	method := requestLineParts[0]
	target := requestLineParts[1]
	versionParts := strings.Split(requestLineParts[2], "/")

	return &RequestLine{
		Method: method, 
		RequestTarget: target,	
		HttpVersion: versionParts[1],	
	}, nil
}
