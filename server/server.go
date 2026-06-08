package server

import (
	"fmt"
	"log"
	"net"
	"strings"
)

const max_buff_size = (4 * 1024)

type Server struct {
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

func NewServer(port string) *Server {
	return &Server{
		port:     port,
		Listener: nil,
	}
}

/*
- Starts the server and begins listening for incoming connections.
- Returns an error if the server fails to start.
*/
func (s *Server) Start() error {
	listener, err := net.Listen("tcp", s.port)
	if err != nil {
		return err
	}
	s.Listener = listener
	log.Printf("Listening on %s", s.port)
	return nil
}

func (s *Server) HandleConnections() {
	if s.Listener == nil {
		log.Printf("Listener not initialized. Call Start() before handling connections.")
		return
	}
	for {
		conn, err := s.Listener.Accept()
		if err != nil {
			log.Printf("Error with connection from %s", conn.RemoteAddr())
			continue
		}
		log.Printf("Handling Connection from %s", conn.RemoteAddr())

		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()
	buff := make([]byte, max_buff_size)
	n, err := conn.Read(buff)
	if err != nil {
		log.Printf("Error reading from connection: %s", err)
		return
	}
	
	req := &Request{Headers: make(map[string]string)}
	parts := strings.SplitN(string(buff[:n]), "\r\n\r\n", 2)
	
	headerLines := strings.Split(parts[0], "\r\n")
	ParseHeader(req, headerLines)
	
	if len(parts) > 1 {
		ParseBody(req, parts[1])
	}
	
	fmt.Printf("Parsed Request: %+v\n", req)
}
