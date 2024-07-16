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
	CONTENT_TYPE_TEXT_PLAIN               = "text/plain"
	CONTENT_TYPE_APPLICATION_OCTET_STREAM = "application/octet-stream"
)

var routes Router = NewRouter()

type Environment struct {
	Directory *string
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	args := os.Args[1:]
	env := Environment{}

	log.Printf("Starting server")
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--directory":
			if i+1 >= len(args) {
				log.Fatalf("No such file or directory")
			}
			env.Directory = &args[i+1]
			log.Printf("Directory: %s", *env.Directory)
		}
	}

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	routes.addHandler(GET, "/", func(req Request, data Data, res *Response) {
		res.Status = 200
	})

	routes.addHandler(GET, "/echo/:text", func(req Request, data Data, res *Response) {
		fmt.Println(req.String())
		res.Status = 200
		res.Headers["Content-Type"] = CONTENT_TYPE_TEXT_PLAIN
		res.Body = []byte(req.Path[6:])
	})

	routes.addHandler(GET, "/files/:filename", func(req Request, data Data, res *Response) {
		if env.Directory != nil {
			filePath := *env.Directory + data["filename"]

			data, err := os.ReadFile(filePath)
			if err != nil {
				log.Printf("Cannot read file: %+v", err)
				res.Status = 500
				return
			}

			res.Status = 200
			res.Headers["Content-Type"] = CONTENT_TYPE_APPLICATION_OCTET_STREAM
			res.Body = data
		}
	})

	routes.addHandler(GET, "/user-agent", func(req Request, data Data, res *Response) {
		res.Status = 200
		res.Headers["Content-Type"] = CONTENT_TYPE_TEXT_PLAIN
		res.Body = []byte(req.Headers["user-agent"])
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
