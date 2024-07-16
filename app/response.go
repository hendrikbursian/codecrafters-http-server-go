package main

import (
	"bytes"
	"fmt"
)

type Response struct {
	Status  httpStatus
	Headers map[string]string
	Body    []byte
}

func (r *Response) Bytes() []byte {
	var out bytes.Buffer

	out.WriteString(fmt.Sprintf("%s %d %s\r\n", HTTP_VERSION, r.Status, httpStauses[httpStatus(r.Status)]))
	for k, v := range r.Headers {
		out.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}
	out.WriteString(fmt.Sprintf("Content-Length: %d\r\n", len(r.Body)))
	out.WriteString("\r\n")
	out.Write(r.Body)

	return out.Bytes()
}

func (r *Response) String() string {
	var out bytes.Buffer
	out.WriteString(fmt.Sprintf("%s %d %s\r\n", HTTP_VERSION, r.Status, httpStauses[r.Status]))
	for k, v := range r.Headers {
		out.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}
	out.WriteString("\r\n")
	out.Write(r.Body)

	return out.String()
}
