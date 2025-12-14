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

func NewWriter(w io.Writer) *Writer {
	return &Writer{
		Writer:      w,
		WriterState: WriteStatusLine,
	}
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	if w.WriterState != WriteStatusLine {
		return fmt.Errorf("error: writing status line in incorrect state: state %v", w.WriterState)
	}

	defer func() {
		w.WriterState = WriteHeaders
	}()

	statusLine := getStatusLine(statusCode)

	_, err := w.Writer.Write(statusLine)

	return err
}

func (w *Writer) WriteChunkedBody(p []byte) (int, error) {
	if w.WriterState != WriteBody {
		return 0, fmt.Errorf("error: writing chunked body in incorrect state: state %v", w.WriterState)
	}

	size := fmt.Sprintf("%x\r\n", len(p)) // in hexadecimal
	data := fmt.Sprintf("%s\r\n", p)

	return w.Writer.Write([]byte(size + data))
}

func (w *Writer) WriteChunkedBodyDone() (int, error) {
	chunkedBodyDone := []byte("0\r\n\r\n")

	return w.Writer.Write(chunkedBodyDone)
}

func (w *Writer) WriteBody(p []byte) (int, error) {
	if w.WriterState != WriteBody {
		return 0, fmt.Errorf("error: writing headers in incorrect state: state %v", w.WriterState)
	}

	defer func() {
		w.WriterState = Done
	}()

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
