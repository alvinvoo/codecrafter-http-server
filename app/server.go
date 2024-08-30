package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/codecrafters-io/http-server-starter-go/app/util"
)

const (
	// This is the port that the server will listen on
	PORT = 4221
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Printf("Server running at port %d \n", PORT)

	l, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", PORT))
	if err != nil {
		fmt.Printf("Failed to bind to port %d \n", PORT)
		os.Exit(1)
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

func handleConnection(conn net.Conn) {
	if tcpConn, ok := conn.(*net.TCPConn); ok {
		defer tcpConn.Close()
		tcpConn.SetKeepAlive(true)
		tcpConn.SetKeepAlivePeriod(10)

		reader := bufio.NewReader(tcpConn)
		requestLine, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading request: ", err.Error())
			return
		}

		if strings.HasPrefix(requestLine, "GET") {
			requestTarget := strings.Split(requestLine, " ")[1]

			var response string
			if requestTarget == "/" {
				response = constructResponse(Response{
					status: Status{
						statusCode: 200,
						statusText: "OK",
					},
				})
			} else if strings.HasPrefix(requestTarget, "/echo/") {
				str := extractPathSegment(requestTarget)
				response = constructResponse(Response{
					status: Status{
						statusCode: 200,
						statusText: "OK",
					}, headers: &Headers{
						contentType: "text/plain",
						contentLen:  len(str),
					}, body: str,
				})
			} else {
				response = constructResponse(Response{
					status: Status{
						statusCode: 404,
						statusText: "Not Found",
					},
				})
			}

			util.DebugLog("Response", response)
			_, err := tcpConn.Write([]byte(response))
			if err != nil {
				fmt.Println("Error writing response: ", err.Error())
				os.Exit(1)
			}
		} else {
			fmt.Println("Unsupported method: ", requestLine)
			return
		}
	}
}

func extractPathSegment(url string) string {
	// Split the URL by "/" and get the part after "echo"
	parts := strings.Split(url, "/")
	if len(parts) > 2 {
		return parts[2] // Get the segment after "/echo/"
	}
	return ""
}

type Status struct {
	statusCode int
	statusText string
}

type Headers struct {
	contentType string
	contentLen  int
}

type Response struct {
	status  Status
	headers *Headers
	body    string
}

func constructResponse(response Response) string {
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
