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
		zw := gzip.NewWriter(&buf)
		_, err := zw.Write(r.Body)
		if err != nil {
			log.Printf("Error while writing gzipped body: %+v", err)
			goto skipzip
		}

		if err := zw.Close(); err != nil {
			log.Printf("Error while closing gzipwriter: %+v", err)
			goto skipzip
		}

		r.Headers["Content-Encoding"] = "gzip"
		r.Headers["Content-Length"] = strconv.Itoa(len(buf.String()))
		r.Body = buf.Bytes()
	} else {
		r.Headers["Content-Length"] = strconv.Itoa(len(string(r.Body)))
	}

skipzip:
	for k, v := range r.Headers {
		out.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}
	out.WriteString("\r\n")

	out.Write(r.Body)

	return out.Bytes()
}
