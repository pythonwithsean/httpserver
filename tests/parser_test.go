package tests

import (
	"testing"

	httpServer "github.com/pythonwithsean/httpserver/server"
)

func TestParseHeader(t *testing.T) {
	header := []string{
		"GET / HTTP/1.1",
		"Host: localhost:5100",
		"User-Agent: curl/7.64.1",
		"Accept: */*",
	}

	req := httpServer.ParseHeader(header)

	if req == nil {
		t.Errorf("Expected non-nil Request object, got nil")
		return
	}

	if req.Method != "GET" {
		t.Errorf("Expected Method 'GET', got '%s'", req.Method)
	}

	if req.Path != "/" {
		t.Errorf("Expected Path '/', got '%s'", req.Path)
	}

	if req.Version != "HTTP/1.1" {
		t.Errorf("Expected Version 'HTTP/1.1', got '%s'", req.Version)
	}

}
