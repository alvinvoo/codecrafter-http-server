package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
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

		// TODO: use Scanner for iteration
		if strings.HasPrefix(requestLine, "GET") {
			requestTarget := strings.Split(requestLine, " ")[1]

			var response string
			if requestTarget == "/" {
				response = "HTTP/1.1 200 OK\r\n\r\n"
			} else {
				response = "HTTP/1.1 404 Not Found\r\n\r\n"
			}
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
