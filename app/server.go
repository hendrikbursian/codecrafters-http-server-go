package main

import (
	"fmt"
	"log"
	"net"
	"os"
)

type httpMethod string

const (
	GET httpMethod = "GET"
)

type httpStatus int

const (
	HTTP_STATUS_NOT_FOUND httpStatus = 404
	HTTP_STATUS_OK        httpStatus = 200
)

var httpStauses = map[httpStatus]string{
	HTTP_STATUS_OK:        "OK",
	HTTP_STATUS_NOT_FOUND: "Not Found",
}

const HTTP_VERSION = "HTTP/1.1"
const (
	CONTENT_TYPE_TEXT_PLAIN = "text/plain"
)

var routes Router

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	routes.handlerOrder = make(map[string]int)
	routes.addHandler(GET, "/", func(req Request, data Data, res *Response) {
		res.Status = 200
	})
	routes.addHandler(GET, "/echo/:text", func(req Request, data Data, res *Response) {
		fmt.Println(req.String())
		res.Status = 200
		res.Headers["Content-Type"] = CONTENT_TYPE_TEXT_PLAIN
		res.Body = req.Path[6:]
	})
	routes.addHandler(GET, "/files", func(req Request, data Data, res *Response) {
		res.Status = 200
	})
	routes.addHandler(GET, "/user-agent", func(req Request, data Data, res *Response) {
		res.Status = 200
		res.Headers["Content-Type"] = CONTENT_TYPE_TEXT_PLAIN
		res.Body = req.Headers["user-agent"]
	})

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
	req := ParseRequest(buf[:n])

	log.Println("Request:\n", req.String())

	res := Response{}
	res.Headers = make(map[string]string)

	handlerFn, data := routes.getHandler(req.Method, req.Path)
	if handlerFn != nil {
		handlerFn(req, data, &res)
	} else {
		res.Status = 404
	}

	log.Println("Response:\n", res.String())

	n, err = conn.Write(res.Bytes())
	if err != nil {
		log.Println("Error sending response: ", err.Error())
	}
}
