package server

import (
	"fmt"
	"io"
	"log"
	"net"
	"sync/atomic"

	"github.com/sithusan/httpfromtcp/internal/request"
	"github.com/sithusan/httpfromtcp/internal/response"
)

type Handler func(w *response.Writer, req *request.Request)

type HandleError struct {
	StatusCode response.StatusCode
	Message    []byte
}

func (hE HandleError) Write(w io.Writer) {
	writer := response.NewWriter(w)
	headers := response.GetDefaultHeaders(len(hE.Message))
	headers.Override("Content-Type", "text/html")

	writer.WriteStatusLine(hE.StatusCode)
	writer.WriteHeaders(headers)
	writer.WriteBody(hE.Message)
}

type Server struct {
	listener net.Listener
	closed   atomic.Bool
	handler  Handler
}

func Serve(port int, handler Handler) (*Server, error) {
	addr := fmt.Sprintf(":%d", port)

	listener, err := net.Listen("tcp", addr)

	if err != nil {
		return nil, err
	}

	server := &Server{
		listener: listener,
		handler:  handler,
	}

	go server.listen()

	return server, nil
}

func (s *Server) Close() error {
	s.closed.Store(true)

	if s.listener != nil {
		return s.listener.Close()
	}

	return nil
}

func (s *Server) listen() {
	for {
		conn, err := s.listener.Accept()

		if err != nil {
			if s.closed.Load() {
				return
			}
			log.Printf("error: listening http %s", err)
			continue
		}

		s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()

	request, err := request.RequestFromReader(conn)

	if err != nil {
		hErr := &HandleError{
			StatusCode: response.BAD_REQUEST,
			Message:    []byte(err.Error()),
		}
		hErr.Write(conn)
		return
	}

	s.handler(
		response.NewWriter(conn),
		request)
}
