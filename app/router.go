package main

import (
	"fmt"
	"strings"
)

type Data map[string]string
type handlerFn func(req Request, data Data, res *Response)

type Router struct {
	handlers     []handlerFn
	handlerOrder map[string]int
}

func NewRouter() Router {
	return Router{
		handlers:     []handlerFn{},
		handlerOrder: make(map[string]int),
	}
}

func (r *Router) getHandlerKey(method string, path string) string {
	return fmt.Sprintf("%s %s", method, path)
}
func (r *Router) addHandler(method string, path string, handler handlerFn) {
	r.handlerOrder[r.getHandlerKey(method, path)] = len(r.handlers)
	r.handlers = append(r.handlers, handler)
}

func validRoute(templatePath []string, path []string) bool {
	if len(path) != len(templatePath) {
		return false
	}

	for i := range templatePath {
		if strings.HasPrefix(templatePath[i], ":") {
			continue
		}
		if templatePath[i] == path[i] {
			continue
		}

		return false
	}

	return true
}

func (r *Router) getHandler(method string, path string) func(req Request, res *Response) {
	data := make(map[string]string)
	pathSplit := strings.Split(path, "/")

	for handlerKey, idx := range r.handlerOrder {
		handlerPath, hasPrefix := strings.CutPrefix(handlerKey, string(method)+" ")
		if !hasPrefix {
			continue
		}

		// check validity
		handlerPathSplit := strings.Split(handlerPath, "/")
		if !validRoute(handlerPathSplit, pathSplit) {
			continue
		}

		// extract data from path /echo/:text -> {text: "..."}
		for i := range len(handlerPathSplit) {
			if key, isDataKey := strings.CutPrefix(handlerPathSplit[i], ":"); isDataKey {
				data[key] = pathSplit[i]
			}
		}

		return func(req Request, res *Response) {
			r.handlers[idx](req, data, res)
		}
	}

	return nil
}
