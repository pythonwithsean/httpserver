package server

import (
	"encoding/json"
	"fmt"
	"net"
	"strings"
	"time"
)

const max_header_size = 8192 // 8KB
const max_chunk_size = 4096  // 4KB
const fifteen_seconds = 15 * time.Second
const thirty_seconds = 30 * time.Second

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

func (s *Server) Start() {
	addr := s.addr + s.port
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		panic(fmt.Sprintf("Error starting server on %s: %v", addr, err))
	}
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
		fmt.Printf("Handling Connection from %s\n", conn.RemoteAddr().String())
		// Set a deadline for the connection to avoid hanging connections
		conn.SetReadDeadline(time.Now().Add(thirty_seconds))
		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()
	var data []byte
	chunk := make([]byte, max_chunk_size) // 4KB buffer
	var headerBlock []byte
	for {
		n, err := conn.Read(chunk)
		if n > 0 {
			data = append(data, chunk[:n]...)
		}
		if err != nil {
			// Check if the error is a timeout error
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				conn.SetWriteDeadline(time.Now().Add(fifteen_seconds))
				conn.Write([]byte("HTTP/1.1 408 Request Timeout\r\n\r\n"))
				fmt.Printf("Timeout reading from connection: %s\n", err)
				return
			}
			// For other errors, send a 400 Bad Request response
			conn.SetWriteDeadline(time.Now().Add(fifteen_seconds))
			conn.Write([]byte("HTTP/1.1 400 Bad Request\r\n\r\n"))
			fmt.Printf("Error reading from connection: %s\n", err)
			return
		}
		// Check if the accumulated data exceeds the maximum header size
		if len(data) > max_header_size {
			conn.SetWriteDeadline(time.Now().Add(fifteen_seconds))
			conn.Write([]byte("HTTP/1.1 413 Payload Too Large\r\n\r\n"))
			fmt.Printf("Header too large from %s\n", conn.RemoteAddr().String())
			return
		}
		// Check if we have received the end of the header section
		if idx := strings.Index(string(data), "\r\n\r\n"); idx != -1 {
			headerBlock = data[:idx]
			break
		}
	}

	if len(headerBlock) == 0 {
		conn.SetWriteDeadline(time.Now().Add(fifteen_seconds))
		conn.Write([]byte("HTTP/1.1 400 Bad Request\r\n\r\n"))
		fmt.Printf("No header received from %s\n", conn.RemoteAddr().String())
		return
	}

	//TODO: reserach CRLF injection and implement a check for it
	req := &Request{Headers: make(map[string]string)}
	ParseHeader(req, strings.Split(string(headerBlock), "\r\n"))
	http_header_json, err := json.MarshalIndent(req.Headers, "", " ")
	if err != nil {
		conn.SetWriteDeadline(time.Now().Add(fifteen_seconds))
		conn.Write([]byte("HTTP/1.1 500 Internal Server Error\r\n\r\n"))
		return
	}
	fmt.Printf("Parsed request from %s:\nMethod: %s\nPath: %s\nVersion: %s\nHost: %s\nHeaders: %s\n", conn.RemoteAddr().String(), req.Method, req.Path, req.Version, req.Host, string(http_header_json))
	conn.SetWriteDeadline(time.Now().Add(fifteen_seconds))
	conn.Write([]byte("HTTP/1.1 200 OK\r\nContent-Length: 0\r\n\r\n"))
}
