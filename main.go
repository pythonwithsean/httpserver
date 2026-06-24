package main

import (
	server "github.com/pythonwithsean/httpserver/server"
)

func main() {
	s := server.NewServer("localhost", ":8000")
	s.Start()
}
