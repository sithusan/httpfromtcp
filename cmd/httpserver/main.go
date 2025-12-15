package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
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
	requestTarget := req.RequestLine.RequestTarget

	if requestTarget == "/yourproblem" {
		writeResponse(w, response.BAD_REQUEST, badRequestResponse)
		return
	}

	if requestTarget == "/myproblem" {
		writeResponse(w, response.INTERNAL_SERVER_ERROR, internalServerResponse)
		return
	}

	httpBinPrefix := "/httpbin"
	if strings.HasPrefix(requestTarget, httpBinPrefix) {
		url := strings.TrimPrefix(requestTarget, httpBinPrefix)
		writeProxyResponse(w, url)
		return
	}

	writeResponse(w, response.OK, successResponse)
}

func writeResponse(w *response.Writer, statusCode response.StatusCode, body []byte) {
	headers := response.GetDefaultHeaders(len(body))
	headers.Override("Content-Type", "text/html")

	w.WriteStatusLine(statusCode)
	w.WriteHeaders(headers)
	w.WriteBody(body)
}

func writeProxyResponse(w *response.Writer, url string) {
	httpBinUrl := fmt.Sprintf("https://httpbin.org%s", url)
	resp, err := http.Get(httpBinUrl)

	if err != nil {
		hErr := &server.HandleError{
			StatusCode: response.INTERNAL_SERVER_ERROR,
			Message:    []byte(err.Error()),
		}
		hErr.Write(w.Writer)
	}

	defer resp.Body.Close()

	buffer := make([]byte, 1024)

	headers := response.GetDefaultHeaders(0)
	headers.Delete("Content-Length")
	headers.Add("Transfer-Encoding", "chunked")

	w.WriteStatusLine(response.OK)
	w.WriteHeaders(headers)
	for {
		n, err := resp.Body.Read(buffer)

		if n > 0 {
			if _, err := w.WriteChunkedBody(buffer[:n]); err != nil {
				fmt.Println("error writing chunked body:", err)
			}
		}

		if err == io.EOF {
			if _, err := w.WriteChunkedBodyDone(); err != nil {
				fmt.Println("error writing body done:", err)
			}
			break
		}

		if err != nil {
			fmt.Printf("error reading body: %s \n", err)
		}
	}
}
