package http

import (
	"fmt"
	"strings"
)

type Status struct {
	StatusCode int
	StatusText string
}

type RespHeaders struct {
	ContentType string
	ContentLen  int
}

type Response struct {
	Status  Status
	Headers *RespHeaders
	Body    string
}

func ExtractPathSegment(url string) string {
	// Split the URL by "/" and get the part after "echo"
	parts := strings.Split(url, "/")
	if len(parts) > 2 {
		return parts[2] // Get the segment after "/echo/"
	}
	return ""
}

func SerializeResponse(response Response) string {
	statusLine := fmt.Sprintf("HTTP/1.1 %d %s\r\n", response.Status.StatusCode, response.Status.StatusText)

	var (
		headers string
		body    string
	)
	if response.Headers != nil {
		headers = fmt.Sprintf("Content-Type: %s\r\nContent-Length: %d\r\n\r\n", response.Headers.ContentType, response.Headers.ContentLen)
	}
	if response.Body != "" {
		body = response.Body
	}

	if response.Headers == nil && response.Body == "" {
		statusLine += "\r\n"
	}

	return fmt.Sprintf("%s%s%s", statusLine, headers, body)
}

func NewNotFoundResponse() Response {
	return Response{
		Status: Status{
			StatusCode: NOT_FOUND,
			StatusText: "Not Found",
		},
	}
}

func NewBadRequestResponse(message string) Response {
	return Response{
		Status: Status{
			StatusCode: BAD_REQUEST,
			StatusText: "Bad Request",
		},
		Headers: &RespHeaders{
			ContentType: "text/plain",
			ContentLen:  len(message),
		},
		Body: message,
	}
}

func NewCreatedResponse() Response {
	return Response{
		Status: Status{
			StatusCode: CREATED,
			StatusText: "Created",
		},
	}
}

func NewSuccessResponse(contentType string, body string) Response {
	return Response{
		Status: Status{
			StatusCode: OK,
			StatusText: "OK",
		},
		Headers: &RespHeaders{
			ContentType: contentType,
			ContentLen:  len(body),
		},
		Body: body,
	}
}
