package response

import (
	"fmt"
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

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	if w.WriterState != WriteHeaders {
		return fmt.Errorf("error: writing headers in incorrect state: state %v", w.WriterState)
	}

	defer func() {
		w.WriterState = WriteBody
	}()

	headerString := ""

	for key, value := range headers {
		headerString += fmt.Sprintf("%s: %s \r\n", key, value)
	}

	headerString += "\r\n"

	_, err := w.Writer.Write([]byte(headerString))

	w.WriterState = WriteBody

	return err
}
