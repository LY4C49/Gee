package gee

import (
	"net/http"
)

// Any object which implements the ServerHTTP implements the http.Handler interface
// It can be used as the second params in http.ListenAndServe

// gin(gee) 框架的核心文件，实现了基本的根据请求的path执行不同的handler功能
// 核心思想是：任何实现了ServerHTTP方法的接口都可以传入给 http.ListenAndServe

/*
ListenAndServe listens on the TCP network address addr
and then calls Serve with handler to handle requests on incoming connections.
*/

/*注释的代码是最原始形态，可以在没有 Router.go 和 context.go 的情况下运行*/

// type HandlerFunc func(w http.ResponseWriter, r *http.Request)
type HandlerFunc func(c *Context)

type Engine struct {
	//router map[string]HandlerFunc
	router *Router
	*RouterGroup
	groups []*RouterGroup
}

func New() *Engine {
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engin: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

func (engine *Engine) addRouter(req_method string, path string, f HandlerFunc) {
	/*
		key := req_method + "-" + path
		engine.router[key] = f
	*/
	engine.router.addRoute(req_method, path, f)
}

func (engine *Engine) GET(path string, f HandlerFunc) {
	engine.addRouter("GET", path, f)
}

func (engine *Engine) POST(path string, f HandlerFunc) {
	engine.addRouter("POST", path, f)
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	/*
		key := r.Method + "-" + r.URL.Path
		if handler, ok := engine.router[key]; ok {
			//c := newContext(w, r)
			//handler(w, r)
		} else {
			fmt.Printf("404 NOT FOUND. Does the path %s really exist?\n", r.URL.Path)
		}
	*/
	c := newContext(w, r)
	engine.router.handle(c)
}

func (engine *Engine) Run(address string) {
	http.ListenAndServe(address, engine)
}
