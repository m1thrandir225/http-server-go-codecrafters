package main

import (
	"fmt"
	"strings"
	// Uncomment this block to pass the first stage
	"net"
	"os"
)

const (
	EOF_MAKER             = "\r\n"
	OK_RESPONSE           = "HTTP/1.1 200 OK\r\n\r\n"
	OK_RESPONSE_WITH_BODY = "HTTP/1.1 200 OK\r\n"
	NOT_FOUND_RESPONSE    = "HTTP/1.1 404 Not Found\r\n\r\n"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")
	// Uncomment this block to pass the first stage
	port := "4221"

	listen, err := net.Listen("tcp", "0.0.0.0:"+port)

	if err != nil {
		fmt.Println("Failed to bind to port " + port)
		os.Exit(1)
	}

	connection, err := listen.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}

	defer listen.Close()

	data := make([]byte, 1024)

	length, err := connection.Read(data)
	if err != nil {
		fmt.Println("Error getting data: ", err.Error())
	}

	requestContent := strings.Split(string(data[:length]), EOF_MAKER)
	reqFirstLine := strings.Split(requestContent[0], " ")

	method := reqFirstLine[0]
	path := reqFirstLine[1]

	var response []byte

	if method == "GET" && path == "/" {
		response = []byte(OK_RESPONSE)
	} else if method == "GET" && strings.HasPrefix(path, "/echo/") {
		expectedString, _ := strings.CutPrefix(path, "/echo/")

		response = []byte(OK_RESPONSE_WITH_BODY + fmt.Sprintf("Content-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s\r\n", len(expectedString), expectedString))
	} else {
		response = []byte(NOT_FOUND_RESPONSE)
	}

	connection.Write(response)
}
