package response

import (
	"fmt"
	"io"
	"strconv"

	"github.com/sithusan/httpfromtcp/internal/headers"
)

type StatusCode int

const (
	OK                    = 200
	BAD_REQUEST           = 400
	INTERNAL_SERVER_ERROR = 500
)

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	switch statusCode {
	case OK:
		_, err := w.Write([]byte("HTTP/1.1 200 OK\r\n"))
		return err
	case BAD_REQUEST:
		_, err := w.Write([]byte("HTTP/1.1 400 Bad Request\r\n"))
		return err
	case INTERNAL_SERVER_ERROR:
		_, err := w.Write([]byte("HTTP/1.1 400 Internal Server Error\r\n"))
		return err
	default:
		_, err := w.Write([]byte("\r\n"))
		return err
	}
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	headers := headers.NewHeaders()
	headers["Content-Length"] = strconv.Itoa(contentLen)
	headers["Content-Type"] = "text/plain"
	headers["Connection"] = "close" // Keep alive will later

	return headers
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	headerString := ""

	for key, value := range headers {
		headerString += fmt.Sprintf("%s: %s \r\n", key, value)
	}

	headerString += "\r\n"

	_, err := w.Write([]byte(headerString))

	return err
}
