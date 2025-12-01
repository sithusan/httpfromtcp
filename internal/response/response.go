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

type writerState int

const (
	WriteStatusLine writerState = iota
	WriteHeaders
	WriteBody
	Done
)

type Writer struct {
	Writer      io.Writer
	WriterState writerState
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	if w.WriterState != WriteStatusLine {
		return fmt.Errorf("error: writing status line in incorrect state: state %v", w.WriterState)
	}

	statusLine := getStatusLine(statusCode)

	_, err := w.Writer.Write(statusLine)

	w.WriterState = WriteHeaders

	return err
}

func (w *Writer) WriteBody(p []byte) (int, error) {

	if w.WriterState != WriteBody {
		return 0, fmt.Errorf("error: writing headers in incorrect state: state %v", w.WriterState)
	}

	w.WriterState = Done

	return w.Writer.Write(p)
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
