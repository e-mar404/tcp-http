package headers

import (
	"bytes"
	"fmt"
	"strings"
)

type Headers map[string]string

func NewHeaders() Headers {
	return Headers{}
}

const crlf = "\r\n"

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	crlfIdx := bytes.Index(data, []byte(crlf))
	if crlfIdx == -1 {
		return 0, false, nil
	}
	
	if crlfIdx == 0 {
		return len(crlf), true, nil
	}

	parts := bytes.SplitN(data[:crlfIdx], []byte(":"), 2)
	fieldName := string(parts[0])
	if fieldName != strings.TrimRight(fieldName, " ") {
		return 0, false, fmt.Errorf("Malformed header, there should be no space between the field key and `:`")
	}

	fieldName = strings.TrimSpace(fieldName)
	fieldValue := strings.TrimSpace(string(parts[1]))

	h.Set(fieldName, fieldValue)

	return crlfIdx + len(crlf), false, nil 
}

func (h Headers) Set(key, value string) {
	h[key] = value
}
