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

	fieldName = strings.ToLower(strings.TrimSpace(fieldName))
	fieldValue := strings.TrimSpace(string(parts[1]))

	err = validate(fieldName)
	if err != nil {
		return 0, false, err 
	}
	h.Set(fieldName, fieldValue)

	return crlfIdx + len(crlf), false, nil 
}

func (h Headers) Set(key, value string) {
	h[key] = value
}

func validate(fieldName string) error {
	specialCharacters := []byte{':', '!', '#', '$', '%', '&', '\'', '*', '+', '-', '.', '^', '_', '`', '|', '~'}

	for _, char := range []byte(fieldName) {
		if !unicode.IsLetter(rune(char)) &&
			!unicode.IsNumber(rune(char)) &&
			!unicode.IsSpace(rune(char)) &&
			bytes.Index(specialCharacters, []byte{char}) == -1 {
				return fmt.Errorf("Invalid character in field name (unicode): %v", char)
		}
	}

	return nil
}
