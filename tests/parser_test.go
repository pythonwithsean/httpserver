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

	req := &httpServer.Request{Headers: make(map[string]string)}
	httpServer.ParseHeader(req, header)

	if req.Method != "GET" {
		t.Errorf("Expected Method 'GET', got '%s'", req.Method)
	}

	if req.Path != "/" {
		t.Errorf("Expected Path '/', got '%s'", req.Path)
	}

	if req.Version != "HTTP/1.1" {
		t.Errorf("Expected Version 'HTTP/1.1', got '%s'", req.Version)
	}

	if req.Host != "localhost:5100" {
		t.Errorf("Expected Host 'localhost:5100', got '%s'", req.Host)
	}

}

func TestIsValidMethod(t *testing.T) {
	tests := []struct {
		method string
		want   bool
	}{
		{"GET", true},
		{"POST", true},
		{"PUT", true},
		{"DELETE", true},
		{"PATCH", true},
		{"HEAD", true},
		{"OPTIONS", true},
		{"INVALID", false},
		{"get", false},
		{"", false},
	}
	for _, tt := range tests {
		got := httpServer.IsValidMethod(tt.method)
		if got != tt.want {
			t.Errorf("IsValidMethod(%q) = %v, want %v", tt.method, got, tt.want)
		}
	}
}

func TestIsValidPath(t *testing.T) {
	tests := []struct {
		path string
		want bool
	}{
		{"/", true},
		{"/users", true},
		{"/users/1", true},
		{"/api/v1/data", true},
		{"users", false},
		{"", false},
	}
	for _, tt := range tests {
		got := httpServer.IsValidPath(tt.path)
		if got != tt.want {
			t.Errorf("IsValidPath(%q) = %v, want %v", tt.path, got, tt.want)
		}
	}
}

func TestIsValidVersion(t *testing.T) {
	tests := []struct {
		version string
		want    bool
	}{
		{"HTTP/1.0", true},
		{"HTTP/1.1", true},
		{"HTTP/2", true},
		{"HTTP/2.0", true},
		{"HTTP/3", false},
		{"http/1.1", false},
		{"", false},
	}
	for _, tt := range tests {
		got := httpServer.IsValidVersion(tt.version)
		if got != tt.want {
			t.Errorf("IsValidVersion(%q) = %v, want %v", tt.version, got, tt.want)
		}
	}
}
