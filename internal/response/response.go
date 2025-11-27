package response

import (
	"fmt"
	"io"
)

type StatusCode int

const (
	OK                    = 200
	BAD_REQUEST           = 400
	INTERNAL_SERVER_ERROR = 500
)

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	statusLine := getStatusLine(statusCode)
	_, err := w.Write(statusLine)

	return err
}

func getStatusLine(statusCode StatusCode) []byte {
	reasonPhrase := ""

	switch statusCode {
	case OK:
		reasonPhrase = "OK"
	case BAD_REQUEST:
		reasonPhrase = "Bad Request"
	case INTERNAL_SERVER_ERROR:
		reasonPhrase = "Internal Server Error"
	}

	return []byte(fmt.Sprintf("HTTP/1.1 %d %s\r\n", statusCode, reasonPhrase))
}
