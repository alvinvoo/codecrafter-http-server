package main

import (
	"bufio"
	"fmt"
	"strings"
)

type RequestLine struct {
	method      string
	url         string
	httpVersion string
}

type ReqHeaders struct {
	host      string
	userAgent string
	accept    string
}

type Request struct {
	requestLine RequestLine
	headers     *ReqHeaders
	body        string
}

func unserializeRequest(reader *bufio.Reader) (Request, error) {
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

	if strings.HasPrefix(requestLine, "GET") {
		fmt.Println("GET request")
		var request Request

		request.requestLine.method = "GET"
		request.requestLine.url = strings.Split(requestLine, " ")[1]
		request.requestLine.httpVersion = strings.Split(requestLine, " ")[2]

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
				break // break straight; there's no keep-alive from client
			}

			nameField := strings.Split(line, ":")[0]
			valueField := strings.TrimSpace(strings.Split(line, ":")[1])

			if request.headers == nil {
				request.headers = &ReqHeaders{}
			}

			switch strings.ToLower(nameField) {
			case "host":
				request.headers.host = valueField
			case "user-agent":
				request.headers.userAgent = valueField
			case "accept":
				request.headers.accept = valueField
			}
		}
		return request, nil
	} else {
		return Request{}, fmt.Errorf("unsupported method: %s", requestLine)
	}
}
