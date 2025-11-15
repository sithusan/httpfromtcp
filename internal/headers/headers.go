package headers

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
)

var crlf = []byte("\r\n")

type Headers map[string]string

func NewHeaders() Headers {
	return map[string]string{}
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	idx := bytes.Index(data, crlf)

	if idx == -1 {
		return 0, false, nil
	}

	if idx == 0 {
		// headers are done, consume the CRLF (CLRF is 2 bytes)
		return len(crlf), true, nil
	}

	parts := bytes.SplitN(data[:idx], []byte(":"), 2)

	if err := validateKey(parts[0]); err != nil {
		return 0, false, err
	}

	h.Set(parts[0], parts[1])

	return idx + len(crlf), false, nil
}

/*
According to RFC, Each field line consists of a case-insensitive field name followed by a colon (":"),
optional leading whitespace, the field line value, and optional trailing whitespace.
*/
func (h Headers) Set(key, value []byte) {
	stringKey := string(bytes.TrimSpace(key))
	stringValue := string(bytes.TrimSpace(value))

	h[strings.ToLower(stringKey)] = stringValue
}

/*
* Helpers
 */
func validateKey(key []byte) error {
	stringKey := string(key)
	errMsg := fmt.Errorf("invalid header name: %s", key)

	// (Not allowed by RFC) Host :
	if stringKey != strings.TrimRight(stringKey, " ") {
		return errMsg
	}

	if len(key) < 1 {
		return errMsg
	}

	pattern := `^[A-Za-z0-9!#$%&'*+\-.\^_` + "`" + `|~ ]+$`

	if match, _ := regexp.MatchString(pattern, stringKey); !match {
		return errMsg
	}

	return nil
}
