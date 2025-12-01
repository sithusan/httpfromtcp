package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/sithusan/httpfromtcp/internal/request"
	"github.com/sithusan/httpfromtcp/internal/response"
	"github.com/sithusan/httpfromtcp/internal/server"
)

const port = 42069

var successResponse = []byte(
	`<html>
<head>
    <title>200 OK</title>
  </head>
  <body>
    <h1>Success!</h1>
    <p>Your request was an absolute banger.</p>
  </body>
</html>
`)

var badRequestResponse = []byte(
	`<html>
<head>
    <title>400 Bad Request</title>
  </head>
  <body>
    <h1>Bad Request</h1>
    <p>Your request honestly kinda sucked.</p>
  </body>
</html>
`)

var internalServerResponse = []byte(
	`<html>
<head>
    <title>500 Internal Server Error</title>
  </head>
  <body>
    <h1>Internal Server Error</h1>
    <p>Okay, you know what? This one is on me.</p>
  </body>
</html>
`)

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

func handler(w *response.Writer, req *request.Request) {
	switch req.RequestLine.RequestTarget {
	case "/yourproblem":
		writeResponse(w, response.BAD_REQUEST, badRequestResponse)
	case "/myproblem":
		writeResponse(w, response.INTERNAL_SERVER_ERROR, internalServerResponse)
	default:
		writeResponse(w, response.OK, successResponse)
	}
}

func writeResponse(w *response.Writer, statusCode response.StatusCode, body []byte) {
	headers := response.GetDefaultHeaders(len(body))
	headers.Override("Content-Type", "text/html")

	w.WriteStatusLine(statusCode)
	w.WriteHeaders(headers)
	w.WriteBody(body)
}
