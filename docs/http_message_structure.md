# HTTP Message Structure (HTTP/1.1)
1. start-line CRL
2. zero or more header field-lines (Name: Value) CRLF
3. CRLF (blank line)
4. optional message-body

-   CRLF is “\r\n” (carriage return + line feed). A Windows/HTTP-style newline.
-   Memory aid: “Registered Nurse” for “\r\n”.


Breakdown:

-   start-line: request line or status line
    -   Example (request): POST /users/primeagen HTTP/1.1
-   field-lines: headers (key-value pairs)
    -   Example: Host: google.com
-   blank line: separates headers from body
-   message-body: optional payload
    -   Example: {"name": "TheHTTPagen"}

Both requests and responses use this structure (with different contents).

## Start-lines vs Headers

-   Start-line (e.g., GET /goodies HTTP/1.1) is not a header.
-   Headers are key-value lines like Host, User-Agent, Accept, Content-Type, Content-Length.

# Key RFCs (Used in this project)

-   [RFC 9110: HTTP Semantics](https://datatracker.ietf.org/doc/html/rfc9110)
-   [RFC 9112: HTTP/1.1](https://datatracker.ietf.org/doc/html/rfc9112)

-   Historical set: RFC 7230–7235 (superseded by 9110/9112), and RFC 2818 (HTTPS over TLS, historical context)