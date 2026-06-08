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

func ParseHeader(header []string) *Request {
	req := &Request{}
	for i, part := range header {
		if i == 0 {
			parts := strings.Split(strings.TrimSpace(part), " ")
			if len(parts) != 3 {
				log.Printf("Invalid request line: expected 3 parts, got %d: %q", len(parts), part)
				return nil
			}
			if !IsValidMethod(parts[0]) {
				log.Printf("Invalid HTTP method: %q", parts[0])
				return nil
			}
			if !IsValidPath(parts[1]) {
				log.Printf("Invalid path: %q (must start with /)", parts[1])
				return nil
			}
			if !IsValidVersion(parts[2]) {
				log.Printf("Invalid HTTP version: %q", parts[2])
				return nil
			}
			req.Method = parts[0]
			req.Path = parts[1]
			req.Version = parts[2]
		}
	}
	return req
}
