package gee

import (
	"log"
	"net/http"
	"path"
)

type RouterGroup struct {
	prefix      string
	middlewares []HandlerFunc
	parent      *RouterGroup
	engin       *Engine
}

// 增加一个 group, 可以理解为创建一个sub-group, root group 是 Engine
func (group *RouterGroup) Group(prefix string) *RouterGroup {
	engine := group.engin
	newGroup := &RouterGroup{
		prefix: group.prefix + prefix,
		parent: group,
		engin:  engine,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

// private method
func (group *RouterGroup) addRoute(method string, path string, f HandlerFunc) {
	full_path := group.prefix + path
	log.Printf("Route %4s - %s", method, full_path)
	group.engin.router.addRoute(method, full_path, f)
}

// === public method ===
func (group *RouterGroup) GET(path string, f HandlerFunc) {
	group.addRoute("GET", path, f)
}

func (group *RouterGroup) POST(path string, f HandlerFunc) {
	group.addRoute("POST", path, f)
}

// middlewares 是基于GroupRoute的，可以给每组添加
func (group *RouterGroup) Use(middlewares ...HandlerFunc) {
	group.middlewares = append(group.middlewares, middlewares...)
}

func (group *RouterGroup) createStaticHandler(relativePath string, fs http.FileSystem) HandlerFunc {
	absolutePath := path.Join(group.prefix, relativePath)
	fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))
	return func(c *Context) {
		file := c.Param("filepath")
		if _, err := fs.Open(file); err != nil {
			c.Status(http.StatusNotFound)
			return
		}
		fileServer.ServeHTTP(c.Writer, c.Req)
	}
}

func (group *RouterGroup) Static(relativePath string, root string) {
	handler := group.createStaticHandler(relativePath, http.Dir(root))
	urlPattern := path.Join(relativePath, "/*filepath")
	// Register GET handlers
	group.GET(urlPattern, handler)
}

/*
	r := gee.New()
	r.Static("/assets", "/usr/geektutu/blog/static")
	// 或相对路径 r.Static("/assets", "./static")
	r.Run(":9999")

	用户访问localhost:9999/assets/js/geektutu.js，
	最终返回/usr/geektutu/blog/static/js/geektutu.js。
*/
