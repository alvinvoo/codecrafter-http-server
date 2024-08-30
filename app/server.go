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
		request, err := unserializeRequest(reader)
		if err != nil {
			fmt.Println("Error reading request: ", err.Error())
			return
		}

		var response string
		if request.requestLine.url == "/" {
			response = serializeResponse(Response{
				status: Status{
					statusCode: 200,
					statusText: "OK",
				},
			})
			// TODO: can extract out the url path segments
		} else if strings.HasPrefix(request.requestLine.url, "/echo") {
			str := extractPathSegment(request.requestLine.url)
			response = serializeResponse(Response{
				status: Status{
					statusCode: 200,
					statusText: "OK",
				}, headers: &RespHeaders{
					contentType: "text/plain",
					contentLen:  len(str),
				}, body: str,
			})
		} else if strings.HasPrefix(request.requestLine.url, "/user-agent") {
			response = serializeResponse(Response{
				status: Status{
					statusCode: 200,
					statusText: "OK",
				}, headers: &RespHeaders{
					contentType: "text/plain",
					contentLen:  len(request.headers.userAgent),
				}, body: request.headers.userAgent,
			})
		} else {
			response = serializeResponse(Response{
				status: Status{
					statusCode: 404,
					statusText: "Not Found",
				},
			})
		}

		util.DebugLog("Response", response)
		_, err = tcpConn.Write([]byte(response))
		if err != nil {
			fmt.Println("Error writing response: ", err.Error())
			os.Exit(1)
		}
	}
}
