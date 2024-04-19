package gee

import "log"

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
