package main

import (
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/sithusan/httpfromtcp/internal/request"
	"github.com/sithusan/httpfromtcp/internal/response"
	"github.com/sithusan/httpfromtcp/internal/server"
)

const port = 42069

func main() {
	server, err := server.Serve(port, handler)

	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()

	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}

// write to the buffer, not the connection.
func handler(w io.Writer, req *request.Request) *server.HandleError {
	requestLine := req.RequestLine

	switch requestLine.RequestTarget {
	case "/yourproblem":
		return &server.HandleError{
			StatusCode: response.BAD_REQUEST,
			Message:    []byte("Your problem is not my problem\n"),
		}
	case "/myproblem":
		return &server.HandleError{
			StatusCode: response.INTERNAL_SERVER_ERROR,
			Message:    []byte("Woopsie, my bad\n"),
		}
	default:
		w.Write(([]byte("All good, frfr\n")))
		return nil
	}
}
