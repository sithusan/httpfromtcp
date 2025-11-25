package server

import (
	"fmt"
	"log"
	"net"
	"sync/atomic"
)

type Server struct {
	listener net.Listener
	closed   atomic.Bool
}

func Serve(port int) (*Server, error) {
	addr := fmt.Sprintf(":%d", port)

	listener, err := net.Listen("tcp", addr)

	if err != nil {
		return nil, err
	}

	server := &Server{
		listener: listener,
	}

	go server.listen()

	return server, nil
}

func (s *Server) Close() error {
	s.closed.Store(true)
	return s.listener.Close()
}

func (s *Server) listen() {
	for {
		conn, err := s.listener.Accept()

		if err != nil {
			log.Fatalf("error: listening http %s", err)
		}

		if s.closed.Load() {
			return
		}

		s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()

	resp := "HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: 13\r\n\r\nHello World!"
	_, _ = conn.Write([]byte(resp))
}
