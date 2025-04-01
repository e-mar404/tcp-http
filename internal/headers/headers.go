package headers

import (
	"bytes"
	"fmt"
	"strings"
	"unicode"
)

type Headers map[string]string

func NewHeaders() Headers {
	return Headers{}
}

const crlf = "\r\n"

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	bytesRead := 0
	sepIdx := bytes.Index(data, []byte(":"))
	if sepIdx == -1 {
		return 0, false, fmt.Errorf("Malformed header, no `:` separating field key and field value")	
	}

	if !unicode.IsLetter(rune(data[sepIdx - 1])) {
		return 0, false, fmt.Errorf("Malformed header, there should be no space between the field key and `:`")
	}

	fieldName := strings.TrimLeft(string(data[:sepIdx]), " ")
	bytesRead += len(fieldName)

	crlfIdx := bytes.Index(data[sepIdx+1:], []byte(crlf))
	if crlfIdx == -1 {
		return 0, false, nil
	}

	fieldValue := strings.Trim(string(data[sepIdx+1:sepIdx+1+crlfIdx]), " ")
	bytesRead = bytesRead + len(fieldValue)

	h[fieldName] = fieldValue

	return bytesRead + 4, false, nil 
}
