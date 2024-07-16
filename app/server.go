package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"slices"
	"strings"
)

var routes Router = NewRouter()
var env Environment = Environment{
	Directory: flag.String("directory", "", ""),
}

type Environment struct {
	Directory *string
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	flag.Parse()

	log.Printf("Starting server")
	if env.Directory != nil {
		log.Printf("Directory: %s", *env.Directory)
	}

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	routes.addHandler("GET", "/", getIndexHandler)
	routes.addHandler("GET", "/echo/:text", getEchoHandler)
	routes.addHandler("GET", "/files/:filename", getFilesHandler)
	routes.addHandler("POST", "/files/:filename", postFilesHandler)
	routes.addHandler("GET", "/user-agent", getUserAgentHandler)

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

	handlerFn := routes.getHandler(req.Method, req.Path)
	if handlerFn != nil {
		handlerFn(req, &res)
	} else {
		res.Status = 404
	}

	acceptedEncodings := strings.Split(req.Headers["accept-encoding"], ", ")
	gzipped := slices.Contains(acceptedEncodings, "gzip")
	if gzipped {
		res.Gzip()
	}

	log.Println("Response:\n", string(res.Bytes()))

	n, err = conn.Write(res.Bytes())
	if err != nil {
		log.Println("Error sending response: ", err.Error())
	}
}

func getIndexHandler(req Request, data Data, res *Response) {
	res.Status = 200
}

func getEchoHandler(req Request, data Data, res *Response) {
	fmt.Println(req.String())
	res.Status = 200
	res.Headers["Content-Type"] = "text/plain"
	res.Body = []byte(data["text"])
}

func getFilesHandler(req Request, data Data, res *Response) {
	if env.Directory == nil {
		res.Status = 404
		return
	}

	filePath := *env.Directory + data["filename"]
	content, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			res.Status = 404
			return
		}

		log.Printf("Cannot read file: %+v", err)
		res.Status = 500
		return
	}

	res.Status = 200
	res.Headers["Content-Type"] = "application/octet-stream"
	res.Body = content
}

func getUserAgentHandler(req Request, data Data, res *Response) {
	res.Status = 200
	res.Headers["Content-Type"] = "text/plain"
	res.Body = []byte(req.Headers["user-agent"])
}

func postFilesHandler(req Request, data Data, res *Response) {
	if env.Directory == nil {
		res.Status = 500
		log.Println("Directory parameter not set!")
		return
	}

	filepath := *env.Directory + data["filename"]

	err := os.WriteFile(filepath, req.Body, 0644)
	if err != nil {
		log.Printf("Error writing file: %+v", err)
		res.Status = 500
		return
	}

	res.Status = 201
}
