package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
)

const (
	EOF_MAKER             = "\r\n"
	OK_RESPONSE           = "HTTP/1.1 200 OK\r\n\r\n"
	OK_RESPONSE_WITH_BODY = "HTTP/1.1 200 OK\r\n"
	NOT_FOUND_RESPONSE    = "HTTP/1.1 404 Not Found\r\n\r\n"
)

func HandleConnection(conn net.Conn, dir string) {
	defer conn.Close()
	req := ReadRequest(conn)

	path, msg := ParsePath(req)

	if path == "/" {
		//200 OK status
		conn.Write([]byte(OK_RESPONSE))
	} else if strings.Contains(path, "echo") {

		response := []byte(OK_RESPONSE_WITH_BODY + fmt.Sprintf("Content-Type: text/plain\r\nContent-length: %d\r\n\r\n%s\r\n", len(msg), msg))

		conn.Write(response)

	} else if strings.Contains(path, "files") {
		file, err := os.ReadFile(dir + msg)
		if err != nil {
			fmt.Printf("Error reading file", err.Error())
			conn.Write([]byte(NOT_FOUND_RESPONSE))
			os.Exit(1)
		}

		fileContent := string(file)

		response := []byte(OK_RESPONSE_WITH_BODY + fmt.Sprintf("Content-Type: application/octet-stream\r\nContent-length: %d\r\n\r\n%s\r\n", len(fileContent), fileContent))

		conn.Write(response)
	} else if strings.Contains(path, "user-agent") {
		userAgent := ParseUserAgent(req)
		response := []byte(OK_RESPONSE_WITH_BODY + fmt.Sprintf("Content-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s\r\n", len(userAgent), userAgent))
		conn.Write(response)
	} else {
		//404 not found
		conn.Write([]byte(NOT_FOUND_RESPONSE))
	}
}

func ReadRequest(conn net.Conn) string {
	data := make([]byte, 1024)

	_, err := conn.Read(data)
	if err != nil {
		fmt.Println("Error getting data: ", err.Error())
	}

	return string(data)
}

/**
* Parse the path
 */
func ParsePath(pathReq string) (string, string) {
	parts := strings.Split(pathReq, " ")

	if len(parts) < 2 {
		fmt.Println("No path to extract")
		os.Exit(1)
	}

	path := parts[1]

	return GetMessageFromPath(path)
}

/*
*Get the message from a given path
 */
func GetMessageFromPath(path string) (string, string) {
	parts := strings.Split(path, "/")
	if len(parts) < 3 {
		return ("/" + parts[1]), ""
	}
	return ("/" + parts[1]), strings.Join(parts[2:], "/")
}

/*
* Return the user agent from an http request
 */
func ParseUserAgent(s string) string {
	parts := strings.Split(s, "\n")

	if len(parts) < 3 {
		fmt.Println("Inorrrect HTTP Request")
		os.Exit(1)
	}

	for _, part := range parts {
		if !strings.Contains(part, "User-Agent") {
			continue
		}

		userAgent := strings.Split(part, ":")
		return strings.TrimSpace(userAgent[1])
	}
	return ""
}

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")
	// Uncomment this block to pass the first stage
	//
	working_dir, err := os.Getwd()

	if err != nil {
		fmt.Printf("Error getting cwd", err.Error())
		os.Exit(1)
	}

	dir := flag.String("directory", working_dir, "files directory")
	flag.Parse()

	port := "4221"

	listen, err := net.Listen("tcp", "0.0.0.0:"+port)

	if err != nil {
		fmt.Println("Failed to bind to port " + port)
		os.Exit(1)
	}

	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		go HandleConnection(conn, *dir)
	}
}
