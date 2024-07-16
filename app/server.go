package main

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"os"
	"slices"
	"strings"
)

const HTTP_VERSION = "HTTP/1.1"
const (
	HTTP_STATUS_NOT_FOUND = "404 Not Found"
	HTTP_STATUS_OK        = "200 OK"
)
const (
	CONTENT_TYPE_TEXT_PLAIN = "text/plain"
)

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Println("Error accepting connection: ", err.Error())
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	var buf []byte = make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		log.Println("Error reading request: ", err.Error())
	}
	req := string(buf[:n])

	log.Println("Request:\n", strings.ReplaceAll(req, "\r\n", "[\\r\\n]"))

	reqSplit := strings.Split(req, "\r\n")

	reqLine := strings.Split(reqSplit[0], " ")
	// verb := reqLine[0]
	path := reqLine[1]
	// version := reqLine[2]

	headerEndIdx := slices.Index(reqSplit, "")
	headers := make(map[string]string)
	for _, h := range reqSplit[1:headerEndIdx] {
		hSplit := strings.Split(h, ": ")
		headers[strings.ToLower(hSplit[0])] = hSplit[1]
	}

	var res bytes.Buffer

	// response line
	res.WriteString(HTTP_VERSION + " ")
	if path == "/" {
		res.WriteString(HTTP_STATUS_OK)
		res.WriteString("\r\n")
		res.WriteString("\r\n")
	} else if strings.HasPrefix(path, "/echo/") {
		body := path[6:]

		res.WriteString(HTTP_STATUS_OK)
		res.WriteString("\r\n")
		res.WriteString(fmt.Sprintf("Content-Type: %s\r\n", CONTENT_TYPE_TEXT_PLAIN))
		res.WriteString(fmt.Sprintf("Content-Length: %d\r\n", len(body)))
		res.WriteString("\r\n")
		res.WriteString(body)
	} else if path == "/user-agent" {
		body := headers["user-agent"]

		res.WriteString(HTTP_STATUS_OK)
		res.WriteString("\r\n")
		res.WriteString(fmt.Sprintf("Content-Type: %s\r\n", CONTENT_TYPE_TEXT_PLAIN))
		res.WriteString(fmt.Sprintf("Content-Length: %d\r\n", len(body)))
		res.WriteString("\r\n")
		res.WriteString(body)
	} else {
		res.WriteString(HTTP_STATUS_NOT_FOUND)
		res.WriteString("\r\n\r\n")
	}

	log.Println("Response:\n", strings.ReplaceAll(res.String(), "\r\n", "[\\r\\n]"))
	n, err = conn.Write(res.Bytes())
	if err != nil {
		log.Println("Error sending response: ", err.Error())
	}
}
