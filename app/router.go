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

func (r *Router) getHandlerKey(method httpMethod, path string) string {
	return fmt.Sprintf("%s %s", method, path)
}
func (r *Router) addHandler(method httpMethod, path string, handler handlerFn) {
	r.handlerOrder[r.getHandlerKey(method, path)] = len(r.handlers)
	r.handlers = append(r.handlers, handler)
}

func (r *Router) getHandler(method httpMethod, path string) (handler handlerFn, routeData map[string]string) {
	data := make(map[string]string)
	pathSplit := strings.Split(path, "/")

handlerLoop:
	for handlerKey, idx := range r.handlerOrder {
		handlerPath, hasPrefix := strings.CutPrefix(handlerKey, string(method)+" ")
		if !hasPrefix {
			continue
		}

		handlerPathSplit := strings.Split(handlerPath, "/")
		if len(pathSplit) != len(handlerPathSplit) {
			continue
		}

		// check validity
		for i := range handlerPathSplit {
			if strings.HasPrefix(handlerPathSplit[i], ":") {
				continue
			}
			if handlerPathSplit[i] == pathSplit[i] {
				continue
			}

			continue handlerLoop
		}

		// extract data
		for i := range len(handlerPathSplit) {
			if after, hasPrefix := strings.CutPrefix(handlerPathSplit[i], ":"); hasPrefix {
				data[after] = pathSplit[i]
			}
		}
		return r.handlers[idx], data
	}

	return nil, nil
}
