package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"log"
)

type Response struct {
	Status  httpStatus
	Headers map[string]string
	Body    []byte
}

func (r *Response) Bytes(gzipped bool) []byte {
	var out bytes.Buffer

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
		r.Body = buf.Bytes()
	}

skipzip:
	out.WriteString(fmt.Sprintf("%s %d %s\r\n", HTTP_VERSION, r.Status, httpStauses[httpStatus(r.Status)]))
	for k, v := range r.Headers {
		out.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}
	out.WriteString(fmt.Sprintf("Content-Length: %d\r\n", len(r.Body)))
	out.WriteString("\r\n")

	out.Write(r.Body)

	return out.Bytes()
}
