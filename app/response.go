package main

import (
	"fmt"
	"strings"
)

type Status struct {
	statusCode int
	statusText string
}

type RespHeaders struct {
	contentType string
	contentLen  int
}

type Response struct {
	status  Status
	headers *RespHeaders
	body    string
}

func extractPathSegment(url string) string {
	// Split the URL by "/" and get the part after "echo"
	parts := strings.Split(url, "/")
	if len(parts) > 2 {
		return parts[2] // Get the segment after "/echo/"
	}
	return ""
}

func serializeResponse(response Response) string {
	statusLine := fmt.Sprintf("HTTP/1.1 %d %s\r\n", response.status.statusCode, response.status.statusText)

	var (
		headers string
		body    string
	)
	if response.headers != nil {
		headers = fmt.Sprintf("Content-Type: %s\r\nContent-Length: %d\r\n\r\n", response.headers.contentType, response.headers.contentLen)
	}
	if response.body != "" {
		body = response.body
	}

	if response.headers == nil && response.body == "" {
		statusLine += "\r\n"
	}

	return fmt.Sprintf("%s%s%s", statusLine, headers, body)
}
