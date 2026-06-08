# Building an HTTP Server from Scratch in Go

> Reference project: **CerberusHTTP** (C++20) by Ronald Adesida
> Goal: Understand what CerberusHTTP does at the OS level, then build the same thing in Go.

---

## What CerberusHTTP Does (Summary)

CerberusHTTP is a C++ HTTP/1.1 server that:
1. Opens a TCP socket on port 8000
2. Uses Linux `epoll` (edge-triggered) to wait for incoming connections
3. Accepts connections, reads raw bytes from each client
4. Parses the HTTP request (start line, headers, JSON body) manually
5. Sends back a hardcoded `HTTP/1.1 200 OK` response
6. Tracks request count with Prometheus metrics

Your Go server will do the exact same thing, but Go's runtime handles the hard OS-level stuff (epoll/kqueue, non-blocking I/O, goroutines) for you. You still need to understand what's happening underneath.

---

## The 3 Pieces

### Piece 1: Sockets & File Descriptors

**What this means:**
A socket is just a file descriptor (an integer) that the OS gives you to represent a network connection. Every `read()`, `write()`, `accept()` call uses this integer to tell the kernel which connection you're talking about.

**What CerberusHTTP does (C++):**
```
getaddrinfo()  →  resolve address
socket()         →  create a file descriptor (returns an int like 3)
bind()           →  attach it to port 8000
listen()         →  mark it as "ready to accept connections" (backlog=10)
accept()         →  when a client connects, get a NEW file descriptor for that connection
```
See: `src/server/tcp.cpp` — `findServerAddress()` and the start of `listenForConnections()`

**What you'll do in Go:**
```go
listener, err := net.Listen("tcp", ":8000")
// Go calls socket(), bind(), listen() internally

conn, err := listener.Accept()
// Go calls accept() internally, returns a net.Conn
```

**What to learn:**
- [ ] What a file descriptor is (open files, sockets, pipes — all fds in Unix)
- [ ] The socket lifecycle: `socket()` → `bind()` → `listen()` → `accept()` → `read()`/`write()` → `close()`
- [ ] Why `bind()` can fail with "address already in use" (and what `SO_REUSEADDR` does)
- [ ] The difference between the **listening socket** (one, stays open) and **connection sockets** (one per client)
- [ ] What `backlog` means in `listen()` — the queue of pending connections

**Build checkpoint:** Write a Go program that listens on `:8000`, accepts a connection, and prints "client connected" — nothing else.

---

### Piece 2: System Calls & I/O (Reading/Writing Data)

**What this means:**
Once a client is connected, you need to read the bytes they send and write bytes back. These are system calls — your program asks the kernel to move data between the network card and your program's memory.

**What CerberusHTTP does (C++):**
```
recv(fd, buffer, 1024, 0)   →  read up to 1024 bytes from the client
send(fd, response, len, 0)   →  write bytes back to the client
```
It uses **non-blocking I/O** (`fcntl` with `O_NONBLOCK`) so `recv()` returns immediately with `EAGAIN` if no data is ready, instead of hanging.

It uses **epoll** to watch many file descriptors at once — when any of them has data ready, `epoll_wait()` returns and tells you which ones.

See: `src/server/tcp.cpp` — `readData()` and `sendResponse()`

**What you'll do in Go:**
```go
buf := make([]byte, 1024)
n, err := conn.Read(buf)    // Go calls read() syscall, blocks until data arrives
data := buf[:n]

conn.Write([]byte("HTTP/1.1 200 OK\r\n..."))  // Go calls write() syscall
```

Go's runtime runs an **epoll/kqueue loop under the hood** (called `netpoll`). When you call `conn.Read()`, Go parks the goroutine and switches to other work. When the kernel says data is ready, Go wakes your goroutine back up. You write blocking-style code but get non-blocking performance.

**What to learn:**
- [ ] What a system call is (user program → kernel boundary)
- [ ] `read()` vs `recv()` and `write()` vs `send()` on sockets
- [ ] Blocking vs non-blocking I/O — why `recv()` can hang forever on a blocking socket
- [ ] I/O multiplexing: `select()` → `poll()` → `epoll` (Linux) / `kqueue` (macOS) — how one thread watches many connections
- [ ] Why CerberusHTTP uses edge-triggered epoll (`EPOLLET`) and reads in a loop until `EAGAIN`
- [ ] How Go's goroutine scheduler + netpoll gives you the same result without manual epoll code

**Build checkpoint:** Read raw bytes from the connected client, print them to stdout, and write a hardcoded response back.

---

### Piece 3: HTTP Parsing (Making Sense of the Bytes)

**What this means:**
The raw bytes you read from the socket are an HTTP request. HTTP/1.1 is a text protocol — you parse it by reading lines and splitting strings. The format is:

```
GET /path HTTP/1.1\r\n
Host: localhost:8000\r\n
Content-Type: application/json\r\n
Content-Length: 27\r\n
\r\n
{"key": "value"}
```

**What CerberusHTTP does (C++):**
1. `appendData()` — accumulate bytes into a string (TCP can split one request across multiple reads)
2. `isRequestComplete()` — check for `\r\n\r\n` (end of headers) + verify `Content-Length` bytes received
3. `parseStartLine()` — split `GET /path HTTP/1.1` into method, path, version
4. `parseHeaders()` — split each line on `:`, trim whitespace, store in a map
5. `parseMessageBody()` — parse JSON body using simdjson

See: `src/parser/parser.cpp` — the full `HttpParser` class

**What you'll do in Go:**
```go
// Read the full request
request := ""
for {
    n, _ := conn.Read(buf)
    request += string(buf[:n])
    if strings.Contains(request, "\r\n\r\n") {
        break
    }
}

// Split headers from body
parts := strings.SplitN(request, "\r\n\r\n", 2)
headerSection := parts[0]
body := parts[1]

// Parse start line
lines := strings.Split(headerSection, "\r\n")
startLine := strings.Split(lines[0], " ")
method := startLine[0]    // "GET"
path := startLine[1]      // "/path"
version := startLine[2]   // "HTTP/1.1"

// Parse headers
headers := make(map[string]string)
for _, line := range lines[1:] {
    kv := strings.SplitN(line, ": ", 2)
    headers[kv[0]] = kv[1]
}
```

**What to learn:**
- [ ] The HTTP/1.1 request format: start line, headers, blank line (`\r\n\r\n`), optional body
- [ ] Why you must handle **partial reads** — TCP is a stream, not a message protocol. One `read()` might give you half a request
- [ ] `Content-Length` — how to know when you've received the full body
- [ ] Building the HTTP response format: status line + headers + body
- [ ] Common status codes: 200 OK, 400 Bad Request, 404 Not Found, 500 Internal Server Error

**Build checkpoint:** Parse method, path, headers, and body from the raw bytes. Send back a proper HTTP response that a browser or `curl` can understand.

---

## Final Project Structure

```
http-server/
├── main.go          # Entry point: listen, accept, handle connections
├── server.go        # TCP listener setup, connection loop
├── parser.go        # HTTP request parsing (start line, headers, body)
├── response.go      # HTTP response builder
├── go.mod
└── requirements.md
```

## Test It With

```bash
# Start your server
go run .

# In another terminal
curl -v http://localhost:8000/
curl -v -X POST -H "Content-Type: application/json" -d '{"hello":"world"}' http://localhost:8000/api
```

## Key Insight

CerberusHTTP manually calls `epoll_create1()`, `epoll_ctl()`, `epoll_wait()`, manages non-blocking fds, and maintains a per-connection parser map. In Go, the runtime does all of this for you — each `conn.Read()` in a goroutine is automatically multiplexed by Go's netpoll. Your job is to understand what's happening underneath so you can debug, optimize, and extend it.
