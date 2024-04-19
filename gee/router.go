package gee

import (
	"log"
	"net/http"
	"strings"
)

// 方便对 router 功能实现增强，实现动态路由等

type Router struct {
	handlers map[string]HandlerFunc
	roots    map[string]*node
}

func newRouter() *Router {
	return &Router{
		handlers: make(map[string]HandlerFunc),
		roots:    make(map[string]*node),
	}
}

func parsePath(path string) []string {
	splited := strings.Split(path, "/")
	parts := make([]string, 0)
	for _, s := range splited {
		if s != "" {
			parts = append(parts, s)
			if s[0] == '*' {
				break
			}
		}
	}
	return parts
}

func (r *Router) addRoute(method string, path string, handler HandlerFunc) {
	log.Printf("Route %s - %s", method, path)
	key := method + "-" + path
	r.handlers[key] = handler

	if _, ok := r.roots[method]; !ok {
		r.roots[method] = &node{}
	}
	parts := parsePath(path)
	r.roots[method].insert(path, parts, 0)
}

func (r *Router) getRoute(method string, path string) (*node, map[string]string) {
	searchParts := parsePath(path)
	params := make(map[string]string)

	root, ok := r.roots[method]
	if !ok {
		return nil, nil
	}

	n := root.search(searchParts, 0)

	if n != nil {
		parts := parsePath(n.full_path)
		for i, part := range parts {
			if part[0] == ':' {
				params[part[1:]] = searchParts[i]
			}
			if part[0] == '*' && len(part) > 1 {
				params[part[1:]] = searchParts[i]
				break // *filepath 是匹配文件，only one!
			}

		}
		return n, params
	}
	return nil, nil
}

func (r *Router) handle(c *Context) {
	// 方法调用
	/*
		key := c.Method + "-" + c.Path
		if handler, ok := r.handlers[key]; ok {
			handler(c)
		} else {
			c.String(http.StatusNotFound, "404 not found: %s\n", c.Path)
		}*/

	node, params := r.getRoute(c.Method, c.Path)
	if node != nil {
		key := c.Method + "-" + node.full_path
		c.Params = params
		r.handlers[key](c)
	} else {
		c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
	}

}
