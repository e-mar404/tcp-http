package headers

type Headers map[string]string

func NewHeaders() Headers {
	return Headers{}
}

const crlf = "\r\n"

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	// get key 

	// go to the end of the line crlf	

	// add that to the key

	// keep going until there are two crlf in a row
}
