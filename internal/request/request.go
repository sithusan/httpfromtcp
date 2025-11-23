package request

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/sithusan/httpfromtcp/internal/headers"
)

var ErrContentBiggerThanLength = errors.New("error: body is bigger than the given content length")

type requestStatus int

const (
	initialized requestStatus = iota
	requestStateParsingHeaders
	requestStateParsingBody
	done
)

const KEY_CONTENT_LENGTH = "Content-Length"

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type Request struct {
	RequestStatus requestStatus
	RequestLine   RequestLine
	Headers       headers.Headers
	Body          []byte
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

func (r *Request) requestParsingBody() bool {
	return r.RequestStatus == requestStateParsingBody
}

func (r *Request) parse(data []byte) (int, error) {
	totalByteParsed := 0

	for !r.done() {
		singleByteParsed, err := r.parseSingle(data[totalByteParsed:])

		if err != nil {
			return 0, err
		}

		totalByteParsed += singleByteParsed

		if singleByteParsed == 0 {
			break
		}
	}

	return totalByteParsed, nil
}

/**
* State Machine
**/

func (r *Request) parseSingle(data []byte) (int, error) {

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

	if r.requestParsingHeaders() {
		consumeBytesFromHeader, headerDone, err := r.Headers.Parse(data)

		if err != nil {
			return 0, err
		}

		if headerDone {
			r.RequestStatus = requestStateParsingBody
		}

		return consumeBytesFromHeader, nil
	}

	if r.requestParsingBody() {

		/**
		According to RFC9110 8.6:
		A user agent SHOULD send Content-Length in a request.
		And "should" has a specific meaning in RFCs per RFC2119:
		This word, or the adjective "RECOMMENDED", mean that
		there may exist valid reasons in particular circumstances to ignore a particular item,
		but the full implications must be understood and carefully weighed before choosing a different course.
		So, going to assume that if there is no Content-Length header, there is no body present.
		**/
		contentLength, ok := r.Headers.Get(KEY_CONTENT_LENGTH)

		if !ok || len(contentLength) == 0 {
			r.RequestStatus = done
			return len(data), nil
		}

		contentLengthInt, err := strconv.Atoi(contentLength)

		if err != nil {
			return 0, err
		}

		if contentLengthInt == 0 {
			r.RequestStatus = done
			return len(data), nil
		}

		r.Body = append(r.Body, data...)

		if len(r.Body) > contentLengthInt {
			return len(data), ErrContentBiggerThanLength
		}

		if len(r.Body) == contentLengthInt {
			r.RequestStatus = done
			return len(data), nil
		}

		return len(data), nil
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
			if errors.Is(err, io.EOF) {
				switch request.RequestStatus {
				case requestStateParsingBody:
					return nil, fmt.Errorf("imcomplete request, body less than content length")
				default:
					return nil, fmt.Errorf("incomplete request, in state: %d, read n bytes on EOF: %d", request.RequestStatus, readedBytes)
				}
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
