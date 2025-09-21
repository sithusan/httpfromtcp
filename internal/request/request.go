package request

import (
	"bytes"
	"fmt"
	"io"
	"strings"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

const CLRF = "\r\n"

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
	request, err := io.ReadAll(reader)

	if err != nil {
		return nil, err
	}

	requestLine, err := parseRequestLine(request)

	if err != nil {
		return nil, err
	}

	return &Request{
		RequestLine: *requestLine,
	}, nil
}

func parseRequestLine(request []byte) (*RequestLine, error) {

	idx, err := checkCLRF(request)

	if err != nil {
		return nil, err
	}

	parts, err := getRequestLineParts(request, idx)

	if err != nil {
		return nil, err
	}

	method, err := getMethod(parts)

	if err != nil {
		return nil, err
	}

	requestTarget, err := getRequestTarget(parts)

	if err != nil {
		return nil, err
	}

	httpVersion, err := getHttpVersion(parts)

	if err != nil {
		return nil, err
	}

	return &RequestLine{
		HttpVersion:   httpVersion,
		RequestTarget: requestTarget,
		Method:        method,
	}, nil
}

func checkCLRF(request []byte) (int, error) {

	idx := bytes.Index(request, []byte(CLRF))

	if idx == -1 {
		return 0, fmt.Errorf("could not find CRLF in request-line")
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
