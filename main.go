package main

import (
	server "github.com/pythonwithsean/httpserver/server"
)

const port = ":5100"

func main() {
	s := server.NewServer(port)
	err := s.Start()
	if err != nil {
		panic(err)
	}
	defer s.Listener.Close()
	s.HandleConnections()

}
