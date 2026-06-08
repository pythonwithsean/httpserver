# HTTP Server from Scratch

Building an HTTP server from scratch in Go to understand operating systems, networking, and the HTTP protocol at a fundamental level.

## What This Project Is

A learning project where I'm implementing an HTTP server without using high-level frameworks. Starting from raw TCP sockets and building up to handle HTTP requests and responses.

## What I'm Learning

### Operating Systems & Networking
- **Socket file descriptors** — How the OS manages network connections as file descriptors
- **Listener vs connection FDs** — The difference between the listening socket (accepts connections) and connection sockets (read/write data)
- **TCP connections** — How clients and servers establish connections with ephemeral ports
- **Port ranges** — Well-known (0-1023), registered (1024-49151), and ephemeral (49152-65535)

### HTTP Protocol
- **Request format** — Parsing raw HTTP requests with `\r\n` line endings
- **Headers and body** — Splitting headers from body using `\r\n\r\n`
- **Request methods** — GET, POST, etc.
- **Response format** — Status lines, headers, and body

### Go Concepts
- **Package structure** — `package main` (programs) vs library packages
- **Testing** — Test files in the same package to access unexported functions
- **Error handling** — Proper nil checking before calling methods
- **Defer** — Cleanup only happens if the defer statement is reached

## Architecture

```
main.go          → Entry point, starts server
server/          → Server logic (listener, connection handling)
├── server.go    → Server struct and methods
└── parser.go    → HTTP request parsing
```

## Running

```bash
make run
# or
go run .
```

## Testing

```bash
curl -v http://localhost:5100/
```

## Roadmap

1. ✅ Set up listener socket and accept connections
2. ✅ Read raw bytes from connection
3. ⏳ Parse HTTP request (method, path, headers, body)
4. ⏳ Build response writer
5. ⏳ Implement router
6. ⏳ Handle incomplete reads and keep-alive connections
