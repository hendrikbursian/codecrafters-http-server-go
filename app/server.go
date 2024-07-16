package main

import (
	"bytes"
	"fmt"
	"net"
	"os"
	"strings"
)

const HTTP_VERSION = "HTTP/1.1"
const (
	HTTP_STATUS_NOT_FOUND = "404 Not Found"
	HTTP_STATUS_OK        = "200 OK"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	conn, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}
	defer conn.Close()

	var buf []byte = make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading request: ", err.Error())
		os.Exit(1)
	}
	req := string(buf[:n])
	fmt.Println("Request: ", req)

	reqSplit := strings.Split(req, "\r\n")

	reqLine := strings.Split(reqSplit[0], " ")
	// verb := reqLine[0]
	path := reqLine[1]
	// version := reqLine[2]

	var res bytes.Buffer

	// response line
	res.WriteString(HTTP_VERSION + " ")
	if path != "/" {
		res.WriteString(HTTP_STATUS_NOT_FOUND)
	} else {
		res.WriteString(HTTP_STATUS_OK)
	}
	res.WriteString("\r\n\r\n")

	n, err = conn.Write(res.Bytes())
	if err != nil {
		fmt.Println("Error sending response: ", err)
		os.Exit(1)
	}
}
