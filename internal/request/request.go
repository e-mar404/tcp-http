package request

import (
	"bytes"
	"fmt"
	"io"
	"regexp"
	"strings"
)

type State int

const (
	Initialized State = iota + 1
	Done
)

type Request struct {
	RequestLine RequestLine
	State       State
}

type RequestLine struct {
	Method        string
	RequestTarget string
	HttpVersion   string
}

const crlf = "\r\n"
const bufSize = 8

func RequestFromReader(reader io.Reader) (*Request, error) {
	buf := make([]byte, bufSize)
	requestBuffer := make([]byte, bufSize)
	readToIndex := 0
	req := &Request{
		State: Initialized,
	}

	for req.State != Done {
		bytesRead, err := reader.Read(buf)
		if err != nil {
			if err == io.EOF {
				fmt.Println("finished reading")
				req.State = Done

				break
			}

			return nil, fmt.Errorf("error reading to buffer: %s", err)
		}

		if readToIndex+bytesRead >= len(requestBuffer) {
			temp := make([]byte, len(requestBuffer)*2)
			copy(temp, requestBuffer[:readToIndex])
			requestBuffer = temp
		}
		readToIndex += bytesRead

		copy(requestBuffer[readToIndex-bytesRead:readToIndex], buf[:bytesRead])

		bytesParsed, err := req.parse(requestBuffer)
		if err != nil {
			return nil, fmt.Errorf("error parsing buf: %s", err)
		}
		
		remainingBytes := readToIndex - bytesParsed
		if bytesParsed > 0 {
			newBuf := make([]byte, bufSize)
			copy(newBuf[:remainingBytes], requestBuffer[bytesRead:])
			requestBuffer = newBuf
		}
		readToIndex = remainingBytes
	}

	return req, nil
}

func (r *Request) parse(data []byte) (int, error) {
	if r.State == Done {
		return 0, fmt.Errorf("trying to parse data that is done")
	}

	if r.State == Initialized {
		requestLine, n, err := parseRequestLine(data)
		if err != nil {
			return 0, err
		}

		if n == 0 {
			return 0, nil
		}

		r.RequestLine = *requestLine
		r.State = Done
		return n, nil
	}

	return 0, fmt.Errorf("state is undefined")
}

func parseRequestLine(data []byte) (*RequestLine, int, error) {
	idx := bytes.Index(data, []byte(crlf))
	if idx == -1 {
		return &RequestLine{}, 0, nil
	}

	requestLineText := string(data[:idx])
	requestLine, err := requestLineFromString(requestLineText)
	if err != nil {
		return nil, 0, fmt.Errorf("error getting request line from text: %v\n", err)
	}

	return requestLine, idx + 2, nil 
}

func requestLineFromString(data string) (*RequestLine, error) {
	requestLineParts := strings.Split(data, " ")

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
		Method:        method,
		RequestTarget: target,
		HttpVersion:   versionParts[1],
	}, nil
}
