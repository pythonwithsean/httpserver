package test

import (
	"testing"

	httpServer "github.com/pythonwithsean/httpserver/main"
)

func TestParseHeader(t *testing.T) {

	header := []string{
		"GET / HTTP/1.1",
		"Host: localhost:5100",
		"User-Agent: curl/7.64.1",
		"Accept: */*",
	}

	httpServer.ParseHeader(header)
}
