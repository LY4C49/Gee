package gee

import (
	"log"
	"net/http"
)

// 方便对 router 功能实现增强，实现动态路由等

type Router struct {
	handlers map[string]HandlerFunc
}

func newRouter() *Router {
	return &Router{handlers: make(map[string]HandlerFunc)}
}

func (r *Router) addRoute(method string, path string, handler HandlerFunc) {
	log.Printf("Route %s - %s", method, path)
	key := method + "-" + path
	r.handlers[key] = handler
}

func (r *Router) handle(c *Context) {
	// 方法调用
	key := c.Method + "-" + c.Path
	if handler, ok := r.handlers[key]; ok {
		handler(c)
	} else {
		c.String(http.StatusNotFound, "404 not found: %s\n", c.Path)
	}
}
