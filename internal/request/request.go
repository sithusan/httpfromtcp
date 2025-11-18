package request

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/sithusan/httpfromtcp/internal/headers"
)

type requestStatus int

const (
	initialized requestStatus = iota
	requestStateParsingHeaders
	done
)

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type Request struct {
	RequestStatus requestStatus
	RequestLine   RequestLine
	Headers       headers.Headers
}

func (r *Request) initialized() bool {
	return r.RequestStatus == initialized
}

func (r *Request) done() bool {
	return r.RequestStatus == done
}

func (r *Request) requestParsingHeaders() bool {
	return r.RequestStatus == requestStateParsingHeaders
}

func (r *Request) parse(data []byte) (int, error) {
	if r.initialized() {
		parsedBytes, err := r.parseRequestLine(data)

		if err != nil {
			return 0, err
		}

		if parsedBytes == 0 {
			return 0, nil
		}

		r.RequestStatus = requestStateParsingHeaders

		return parsedBytes, nil
	}

	// Done from parse = real header parsing done.
	// We need to keep calling header parse, until consume byte from header parse is equal to current
	// chunk of the data. Because single chunk can contain multiple headers and we split just by N.
	if r.requestParsingHeaders() {
		totalConsumeBytesFromHeader := 0

		for {
			consumeBytesFromHeader, headerDone, err := r.Headers.Parse(data[totalConsumeBytesFromHeader:])

			if err != nil {
				return 0, err
			}

			if headerDone {
				r.RequestStatus = done
				return totalConsumeBytesFromHeader, nil
			}

			totalConsumeBytesFromHeader += consumeBytesFromHeader

			if consumeBytesFromHeader == 0 {
				// headers are done, needs to return totalConsume bytes, because of the len of CLRF.
				return totalConsumeBytesFromHeader, nil
			}

			if totalConsumeBytesFromHeader == len(data) {
				// Consumed all available data, but headers not done yet
				return totalConsumeBytesFromHeader, nil
			}
		}

	}

	if r.done() {
		return 0, fmt.Errorf("error: trying to read the data in done state")
	}

	return 0, fmt.Errorf("error: unknown state")
}

func (r *Request) parseRequestLine(data []byte) (int, error) {

	idx, err := checkCLRF(data)

	if err != nil {
		return 0, err
	}

	// just needs more data
	if idx == 0 {
		return 0, nil
	}

	parts, err := getRequestLineParts(data, idx)

	if err != nil {
		return 0, err
	}

	method, err := getMethod(parts)

	if err != nil {
		return 0, err
	}

	requestTarget, err := getRequestTarget(parts)

	if err != nil {
		return 0, err
	}

	httpVersion, err := getHttpVersion(parts)

	if err != nil {
		return 0, err

	}

	r.RequestLine = RequestLine{
		HttpVersion:   httpVersion,
		RequestTarget: requestTarget,
		Method:        method,
	}

	return idx + len(CLRF), nil
}

func NewRequest() *Request {
	return &Request{
		RequestStatus: initialized,
		Headers:       headers.NewHeaders(),
	}
}

var CLRF = []byte("\r\n")

const VERSION_PART_SIZE = 2

var requestLineParts = map[string]int{
	"method":        0,
	"requestTarget": 1,
	"httpVersion":   2,
}

// empty struct allocate zero memory
var supportedMethods = map[string]struct{}{
	"GET":    {},
	"POST":   {},
	"PUT":    {},
	"PATCH":  {},
	"DELETE": {},
	"HEAD":   {},
}

func RequestFromReader(reader io.Reader) (*Request, error) {

	buffer := make([]byte, 8)
	readToIndex := 0
	request := NewRequest()

	for !request.done() {
		// buffer resizing
		if len(buffer) <= readToIndex {
			newBuffer := make([]byte, (len(buffer) * 2))
			copy(newBuffer, buffer)
			buffer = newBuffer
		}

		readedBytes, err := reader.Read(buffer[readToIndex:])

		if err != nil {
			if err == io.EOF {
				request.RequestStatus = done
				break
			}
			return nil, err
		}

		readToIndex += readedBytes
		parsedBytes, err := request.parse(buffer[:readToIndex])

		if err != nil {
			return nil, err
		}

		// remove the used ones
		// Shift the unparsed data (from parsedBytes onwards) to the beginning of the buffer
		// Example: if parsedBytes is 3, copy buffer[3:] to buffer[0:]
		copy(buffer, buffer[parsedBytes:readToIndex])
		readToIndex -= parsedBytes
	}

	return request, nil
}

/**
* Helpers
**/

// No error, because that just means that it needs more data before it can parse the request line.
func checkCLRF(request []byte) (int, error) {

	idx := bytes.Index(request, CLRF)

	if idx == -1 {
		return 0, nil
	}

	return idx, nil
}

func getRequestLineParts(request []byte, idx int) ([]string, error) {

	requestLine := string(request[:idx])

	parts := strings.Split(requestLine, " ")

	if len(parts) != len(requestLineParts) {
		return nil, fmt.Errorf("not enough parts")
	}

	return parts, nil
}

func getMethod(parts []string) (string, error) {
	method := parts[requestLineParts["method"]]

	for _, char := range method {
		if char < 'A' || char > 'Z' {
			return "", fmt.Errorf("method needs be in capital")
		}
	}

	if _, ok := supportedMethods[method]; !ok {
		return "", fmt.Errorf("unsupported method")
	}

	return method, nil
}

func getRequestTarget(parts []string) (string, error) {
	requestTarget := parts[requestLineParts["requestTarget"]]

	if requestTarget == "" {
		return "", fmt.Errorf("malformed request target")
	}

	return requestTarget, nil
}

func getHttpVersion(parts []string) (string, error) {
	versionParts := strings.Split(parts[requestLineParts["httpVersion"]], "/")

	if len(versionParts) != VERSION_PART_SIZE {
		return "", fmt.Errorf("malformed http version")
	}

	if versionParts[0] != "HTTP" {
		return "", fmt.Errorf("unrecognized http version")
	}

	// Currently, ONLY SUPPORT 1.1.
	if versionParts[1] != "1.1" {
		return "", fmt.Errorf("unrecognized http version")
	}

	return versionParts[1], nil
}
