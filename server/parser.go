package server

import (
	"log"
	"strings"
)

var validMethods = map[string]bool{
	"GET":     true,
	"POST":    true,
	"PUT":     true,
	"DELETE":  true,
	"PATCH":   true,
	"HEAD":    true,
	"OPTIONS": true,
}

var validVersions = map[string]bool{
	"HTTP/1.0": true,
	"HTTP/1.1": true,
	"HTTP/2":   true,
	"HTTP/2.0": true,
}

func IsValidMethod(method string) bool {
	return validMethods[method]
}

func IsValidPath(path string) bool {
	return strings.HasPrefix(path, "/")
}

func IsValidVersion(version string) bool {
	return validVersions[version]
}

func Empty(s string) bool {
	return strings.TrimSpace(s) == ""
}

func IsValidHeaderKey(key string) bool {
	return !Empty(key) && !strings.ContainsAny(key, ": ")
}

func ParseHeaderField(line string) (string, string, bool) {
	parts := strings.SplitN(line, ":", 2)
	if len(parts) != 2 {
		return "", "", false
	}
	key := strings.TrimSpace(parts[0])
	value := strings.TrimSpace(parts[1])
	if Empty(key) || Empty(value) {
		return "", "", false
	}
	return key, value, true
}

func ParseHeader(req *Request, header []string) {
	if len(header) == 0 {
		return
	}

	parts := strings.Split(strings.TrimSpace(header[0]), " ")
	if len(parts) != 3 {
		log.Printf("Invalid request line: expected 3 parts, got %d: %q", len(parts), header[0])
		return
	}
	if !IsValidMethod(parts[0]) {
		log.Printf("Invalid HTTP method: %q", parts[0])
		return
	}
	if !IsValidPath(parts[1]) {
		log.Printf("Invalid path: %q (must start with /)", parts[1])
		return
	}
	if !IsValidVersion(parts[2]) {
		log.Printf("Invalid HTTP version: %q", parts[2])
		return
	}
	req.Method = parts[0]
	req.Path = parts[1]
	req.Version = parts[2]

	for i := 1; i < len(header); i++ {
		key, value, ok := ParseHeaderField(header[i])
		if !ok {
			log.Printf("Invalid header: %q", header[i])
			continue
		}
		req.Headers[key] = value
		if key == "Host" {
			req.Host = value
		}
	}
}

func ParseBody(req *Request, body string) {
	if !Empty(body) {
		req.Body = body
	}
}
