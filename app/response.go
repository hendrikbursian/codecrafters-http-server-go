package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"log"
)

var httpStatuses = map[int]string{
	200: "OK",
	404: "Not Found",
	201: "Created",
}

type Response struct {
	Status  int
	Headers map[string]string
	Body    []byte
}

func (r *Response) Gzip() {
	var gzBuf bytes.Buffer
	zw := gzip.NewWriter(&gzBuf)

	if _, err := zw.Write(r.Body); err != nil {
		log.Printf("Error while writing to gzip writer: %+v", err)
		return
	}

	if err := zw.Close(); err != nil {
		log.Printf("Error while closing gzip writer: %+v", err)
		return
	}

	r.Body = gzBuf.Bytes()
	r.Headers["Content-Encoding"] = "gzip"
}

func (r *Response) Bytes() []byte {
	var out bytes.Buffer

	out.WriteString(fmt.Sprintf("%s %d %s\r\n", "HTTP/1.1", r.Status, httpStatuses[r.Status]))
	for k, v := range r.Headers {
		out.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}
	out.WriteString(fmt.Sprintf("%s: %d\r\n", "Content-Length", len(r.Body)))
	out.WriteString("\r\n")
	out.Write(r.Body)
	return out.Bytes()
}
