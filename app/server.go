package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/codecrafters-io/http-server-starter-go/app/http"
	"github.com/codecrafters-io/http-server-starter-go/app/util"
)

const (
	// This is the port that the server will listen on
	PORT = 4221
)

var fileDirectory string

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Printf("Server running at port %d \n", PORT)

	l, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", PORT))
	if err != nil {
		fmt.Printf("Failed to bind to port %d \n", PORT)
		os.Exit(1)
	}

	if len(os.Args) > 1 && os.Args[1] == "--directory" {
		fileDirectory = os.Args[2]
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		go handleConnection(conn)
	}
}

func handleEncoding(request *http.Request, response *http.Response) {
	if len(request.Headers.AcceptEncoding) != 0 && response.Body != "" {
		// since response's body has something, response header should alrdy be set
		for _, encoding := range request.Headers.AcceptEncoding {
			if encoding == "gzip" {
				// Compress the response body using gzip
				// and set the Content-Encoding header to "gzip"

				response.Headers.ContentEncoding = append(response.Headers.ContentEncoding, "gzip")
				body, err := util.GzipEncode(response.Body)
				if err != nil {
					fmt.Println("Error compressing response body: ", err.Error())
					response = http.NewInternalServerErrorResponse(err.Error())
				}

				response.Body = body
				response.Headers.ContentLen = len(response.Body)
			}
		}
	}
}

func generateResponse(request *http.Request) *http.Response {
	var response *http.Response

	if request.RequestLine.Url == "/" {
		response = &http.Response{
			Status: http.Status{
				StatusCode: http.OK,
				StatusText: "OK",
			},
		}
	} else if strings.HasPrefix(request.RequestLine.Url, "/echo") {
		str := http.ExtractPathSegment(request.RequestLine.Url)
		response = http.NewSuccessResponse("text/plain", str)
	} else if strings.HasPrefix(request.RequestLine.Url, "/user-agent") {
		response = http.NewSuccessResponse("text/plain", request.Headers.UserAgent)
	} else if strings.HasPrefix(request.RequestLine.Url, "/files") {
		fileName := http.ExtractPathSegment(request.RequestLine.Url)

		if request.RequestLine.Method == "POST" {
			if request.Headers.ContentType != "application/octet-stream" {
				return http.NewBadRequestResponse("Content-Type must be application/octet-stream")
			}

			if request.Headers.ContentLen == 0 || request.Body == "" {
				return http.NewBadRequestResponse("Content-Length must be greater than 0")
			}

			if request.Headers.ContentLen != len(request.Body) {
				return http.NewBadRequestResponse("Content-Length does not match body length")
			}

			err := os.WriteFile(fmt.Sprintf("%s/%s", fileDirectory, fileName), []byte(request.Body), 0644)
			if err != nil {
				fmt.Println("Error writing file: ", err.Error())
				response = http.NewBadRequestResponse(err.Error())
			} else {
				response = http.NewCreatedResponse()
			}
		} else {
			fileContent, err := os.ReadFile(fmt.Sprintf("%s/%s", fileDirectory, fileName))
			if err != nil {
				fmt.Println("Error reading file: ", err.Error())
				response = http.NewNotFoundResponse()
			} else {
				response = http.NewSuccessResponse("application/octet-stream", string(fileContent))
			}
		}
	} else {
		response = http.NewNotFoundResponse()
	}

	handleEncoding(request, response)

	return response
}

func handleConnection(conn net.Conn) {
	if tcpConn, ok := conn.(*net.TCPConn); ok {
		defer tcpConn.Close()
		tcpConn.SetKeepAlive(true)
		tcpConn.SetKeepAlivePeriod(10)

		reader := bufio.NewReader(tcpConn)
		request, err := http.UnserializeRequest(reader)
		if err != nil {
			fmt.Println("Error reading request: ", err.Error())
			return
		}

		serializedResponse := http.SerializeResponse(generateResponse(&request))

		util.DebugLog("Response", serializedResponse)
		_, err = tcpConn.Write([]byte(serializedResponse))
		if err != nil {
			fmt.Println("Error writing response: ", err.Error())
			os.Exit(1)
		}
	}
}
