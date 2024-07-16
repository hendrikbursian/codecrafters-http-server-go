package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"log"
	"strconv"
)

type Response struct {
	Status  httpStatus
	Headers map[string]string
	Body    []byte
}

func (r *Response) Bytes(gzipped bool) []byte {
	var out bytes.Buffer

	out.WriteString(fmt.Sprintf("%s %d %s\r\n", HTTP_VERSION, r.Status, httpStauses[httpStatus(r.Status)]))

	if gzipped {
		var buf bytes.Buffer
		enc := gzip.NewWriter(&buf)
		enc.Write(r.Body)
		enc.Close()
		log.Println(buf.String())

		contentLength := strconv.Itoa(len(buf.String()))
		log.Println(contentLength)

		r.Headers["Content-Length"] = contentLength
		r.Headers["Content-Encoding"] = "gzip"
		r.Body = buf.Bytes()
	} else {
		r.Headers["Content-Length"] = strconv.Itoa(len(string(r.Body)))
	}

	for k, v := range r.Headers {
		out.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}
	out.WriteString("\r\n")

	out.Write(r.Body)

	return out.Bytes()
}
