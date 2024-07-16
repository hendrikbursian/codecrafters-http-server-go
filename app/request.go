package main

import (
	"bytes"
	"fmt"
	"log"
	"slices"
	"strconv"
	"strings"
)

type Request struct {
	Method  httpMethod
	Version string
	Path    string
	Headers map[string]string
	Body    []byte
}

func (r *Request) String() string {
	var out bytes.Buffer
	out.WriteString(fmt.Sprintf("%s %s %s\r\n", r.Method, r.Path, r.Version))
	out.WriteString("Headers:\n")
	for k, v := range r.Headers {
		out.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}
	out.WriteString("\r\n")
	out.Write(r.Body)

	return out.String()
}

func ParseRequest(bytes []byte) Request {
	req := string(bytes)
	reqSplit := strings.Split(req, "\r\n")
	reqLine := strings.Split(reqSplit[0], " ")

	headerEndIdx := slices.Index(reqSplit, "")
	headers := make(map[string]string)
	for _, h := range reqSplit[1:headerEndIdx] {
		hSplit := strings.Split(h, ": ")
		headers[strings.ToLower(hSplit[0])] = hSplit[1]
	}

	ret := Request{
		Method:  httpMethod(reqLine[0]),
		Path:    reqLine[1],
		Version: reqLine[2],
		Headers: headers,
	}

	if contentLength, ok := headers["content-length"]; ok {
		contentLength, err := strconv.Atoi(contentLength)
		if err != nil {
			log.Printf("Cannot parse content-length: %s", err)
			return ret
		}

		ret.Body = []byte(reqSplit[headerEndIdx+1])[:contentLength]
	}

	return ret
}
