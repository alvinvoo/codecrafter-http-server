package http

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
)

type RequestLine struct {
	Method      string
	Url         string
	HttpVersion string
}

type ReqHeaders struct {
	Host           string
	UserAgent      string
	Accept         string
	ContentType    string
	ContentLen     int
	AcceptEncoding string
}

type Request struct {
	RequestLine RequestLine
	Headers     *ReqHeaders
	Body        string
}

func UnserializeRequest(reader *bufio.Reader) (Request, error) {
	// Read the request line (e.g., "GET /path HTTP/1.1\r\n")
	requestLine, err := reader.ReadString('\n')
	if err != nil {
		if err.Error() == "EOF" {
			fmt.Println("Client disconnected")
		} else {
			fmt.Println("Error reading request line: ", err.Error())
		}
		return Request{}, err
	}

	if strings.HasPrefix(requestLine, "GET") || strings.HasPrefix(requestLine, "POST") {
		var request Request

		requestParts := strings.Split(requestLine, " ")
		request.RequestLine.Method = requestParts[0]
		request.RequestLine.Url = requestParts[1]
		request.RequestLine.HttpVersion = requestParts[2]

		for {
			content, err := reader.ReadString('\n')
			if err != nil {
				if err.Error() == "EOF" {
					fmt.Println("Client disconnected")
				} else {
					fmt.Println("Error reading request line: ", err.Error())
				}
				break
			}

			// Check if the request line is properly terminated
			line := strings.TrimSuffix(content, "\r\n")
			if line == "" {
				break // End of headers
			}

			if strings.ContainsRune(line, rune(':')) {
				nameField := strings.Split(line, ":")[0]
				valueField := strings.TrimSpace(strings.Split(line, ":")[1])

				if request.Headers == nil {
					request.Headers = &ReqHeaders{}
				}

				switch strings.ToLower(nameField) {
				case "host":
					request.Headers.Host = valueField
				case "user-agent":
					request.Headers.UserAgent = valueField
				case "accept":
					request.Headers.Accept = valueField
				case "content-type":
					request.Headers.ContentType = valueField
				case "content-length":
					length, err := strconv.Atoi(valueField)
					if err != nil {
						return Request{}, fmt.Errorf("error converting content length to integer: %s", err.Error())
					}
					request.Headers.ContentLen = length
				case "accept-encoding":
					request.Headers.AcceptEncoding = valueField
				}
			}
		}

		if request.Headers != nil && request.Headers.ContentLen > 0 {
			bodyBytes := make([]byte, request.Headers.ContentLen)
			n, err := reader.Read(bodyBytes)
			if err != nil {
				fmt.Println("Error reading request body: ", err.Error())
				return Request{}, err
			}

			if n != request.Headers.ContentLen {
				return Request{}, fmt.Errorf("unexpected content length, got %d bytes, expected %d", n, request.Headers.ContentLen)
			}

			request.Body = string(bodyBytes)
		}
		return request, nil
	} else {
		return Request{}, fmt.Errorf("unsupported method: %s", requestLine)
	}
}
