package server

import (
	"errors"
	"fmt"
	"io"
	"net"
	"strings"
	"time"
)

const max_header_size = 8192 // 8KB
const max_chunk_size = 1     // 1 Byte
const max_conn_duration = 30 * time.Second
const CRLF = "\r\n"

type Server struct {
	addr     string
	port     string
	Listener net.Listener
}

type Request struct {
	Method  string
	Path    string
	Version string
	Host    string
	Headers map[string]string
	Body    string
}

func NewServer(addr, port string) *Server {
	return &Server{
		addr:     addr,
		port:     port,
		Listener: nil,
	}
}

func newListener(addr string) net.Listener {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		panic(fmt.Sprintf("Error starting server on %s: %v", addr, err))
	}
	return listener
}

func (s *Server) Start() {
	addr := s.addr + s.port
	listener := newListener(addr)
	s.Listener = listener
	defer s.Listener.Close()
	fmt.Printf("Listening on %s\n", addr)
	s.HandleConnections()
}

func (s *Server) HandleConnections() {
	if s.Listener == nil {
		panic("Listener is not initialized")
	}
	for {
		conn, err := s.Listener.Accept()
		if err != nil {
			fmt.Printf("Error with connection from %s\n", conn.RemoteAddr().String())
			continue
		}
		// fmt.Printf("Handling Connection from %s\n", conn.RemoteAddr().String())
		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()
	var data []byte
	chunk := make([]byte, max_chunk_size) // 1 byte buffer
	var headerBlock []byte
	var bodyBlockIdx int
	// if json header has connection: keep-alive, we should keep the connection open for a certain duration
	conn.SetDeadline(time.Now().Add(max_conn_duration)) // Set a deadline for the connection
	// TODO: Improve performance of searching for CRLF+CRlF
	for {
		// Set a deadline for the connection to avoid hanging connections
		n, err := conn.Read(chunk)
		if n > 0 {
			data = append(data, chunk[:n]...)
		}
		if err != nil {
			// Check if the error is a timeout error
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				conn.Write([]byte("HTTP/1.1 408 Request Timeout\r\n\r\n"))
				fmt.Printf("Timeout reading from connection: %s\n", err)
				return
			}
			if errors.Is(err, io.EOF) {
				// Client closed the connection; process whatever data we have
				break
			}
			// For other errors, send a 400 Bad Request response
			conn.Write([]byte("HTTP/1.1 400 Bad Request\r\n\r\n"))
			fmt.Printf("Error reading from connection: %s\n", err)
			return
		}
		// Check if the accumulated data exceeds the maximum header size
		if len(data) > max_header_size {
			conn.Write([]byte("HTTP/1.1 413 Payload Too Large\r\n\r\n"))
			fmt.Printf("Header too large from %s\n", conn.RemoteAddr().String())
			return
		}
		// Check if we have received the end of the header section
		if idx := strings.Index(string(data), CRLF+CRLF); idx != -1 {
			headerBlock = data[:idx]
			bodyBlockIdx = idx + len(CRLF+CRLF)
			break
		}
	}

	if len(headerBlock) == 0 {
		conn.Write([]byte("HTTP/1.1 400 Bad Request\r\n\r\n"))
		fmt.Printf("No header received from %s\n", conn.RemoteAddr().String())
		return
	}

	//TODO: reserach CRLF injection and implement a check for it
	req := &Request{Headers: make(map[string]string)}
	ParseHeader(req, strings.Split(string(headerBlock), CRLF))
	if req.Method == "" || req.Path == "" || req.Version == "" || req.Host == "" || len(req.Headers) == 0 {
		conn.Write([]byte("HTTP/1.1 400 Bad Request\r\n\r\n"))
		fmt.Printf("Invalid request from %s\n", conn.RemoteAddr().String())
		return
	}

	fmt.Printf("Content-Length %s, bodyBlockIdx %d\n", req.Headers["content-type"], bodyBlockIdx)
	fmt.Printf("Body: %s\n", string(data[bodyBlockIdx:]))
	// cl, err := strconv.Atoi(req.Headers["content-type"])
	// if err != nil {
	// 	conn.Write([]byte("HTTP/1.1 400 Bad Request\r\n\r\n"))
	// 	fmt.Printf("Error reading from connection: %s\n", err)
	// }

	// fmt.Printf("Parsed request from %s \nMethod: %s\nPath: %s\nVersion: %s\nHost: %s\nHeaders: %s\n", conn.RemoteAddr().String(), req.Method, req.Path, req.Version, req.Host, string(http_header_json))
	body := "<html><body><h1>Hello, Sean!</h1></body></html>"
	conn.Write([]byte(fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Length: %d\r\nContent-Type: text/html\r\n\r\n%s", len(body), body)))
}
