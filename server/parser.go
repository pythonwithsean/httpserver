package server

import (
	"fmt"
	"strings"
)

func ParseHeader(header []string) *Request {
	fmt.Printf("Header Lines: [%v]\n", strings.Join(header, ", "))
	return nil
}
