package response

import (
	"fmt"
	"io"
	"strconv"

	"github.com/sithusan/httpfromtcp/internal/headers"
)

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
