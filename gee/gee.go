package gee

import (
	"html/template"
	"net/http"
	"strings"
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

	htmlTemplates *template.Template
	funcMap       template.FuncMap
}

func New() *Engine {
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engin: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

func Default() *Engine {
	engine := New()
	engine.Use(Logger(), Recovery())
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

	/*
		c := newContext(w, r)
		engine.router.handle(c)
	*/

	/*
		一个请求从ServerHTTP进入
		1.先检测它的URL上是否属于某个GROUP (是否包含GROUP的前缀)
		2.找齐所有满足的GROUP，将该GROUP的middlewares添加到调用链中
		3.ServerHTTP最后调用router.handle
		4.router会通过前缀树查找该路径是否有绑定的方法(HandleFunc)。若有，则append进入context.middlewares调用链中。没有则结束，代表当前URL无绑定方法
		5.router.handle最后会调用c.Next 开始依次执行[]middlewares中的所有HandleFunc
		***在上述一系列过程中，Context承担了居中传递信息的功能，从找齐所有GROUP的middlewares，到传给route.handle, 均利用context(context有一个成员 []middlewares)***
		// 注意：若在HandlerFunc中调用c.Next()则可以实现 “中断” 的效果
	*/
	var middlewares []HandlerFunc
	for _, group := range engine.groups {
		if strings.HasPrefix(r.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}
	c := newContext(w, r)
	c.handlers = middlewares
	c.engine = engine
	engine.router.handle(c)

}

func (engine *Engine) Run(address string) {
	http.ListenAndServe(address, engine)
}

func (engine *Engine) SetFuncMap(funcMap template.FuncMap) {
	engine.funcMap = funcMap
}

func (engine *Engine) LoadHTMLGlob(pattern string) {
	engine.htmlTemplates = template.Must(template.New("").Funcs(engine.funcMap).ParseGlob(pattern))
}
