package gee

import (
	"log"
	"net/http"
)

// HandlerFunc defines the request handler used by gee
type HandlerFunc func(*Context)

// Engine implement the interface of ServeHTTP
type (
	Engine struct {
		*RouterGroup
		router *router
		groups []*RouterGroup
	}

	RouterGroup struct {
		prefix      string
		middlewares []HandlerFunc
		parent      *RouterGroup
		engine      *Engine
	}
)

// New is the constructor of gee.Engine
func New() *Engine {
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}

	return engine
}

func (group *RouterGroup) Group(prefix string) *RouterGroup {
	engine := group.engine
	newGroup := &RouterGroup{
		prefix: group.prefix + prefix,
		parent: group,
		engine: engine,
	}
	engine.groups = append(engine.groups, newGroup)

	return newGroup
}

// 注意, RouterGroup 的作用只是管理 group 和衔接 engine，
// 而避免了同类型路由的冗余
// Router 访问真正的执行还是需要通过 group.engine 来实现
// 这也就是添加 RouterGroup 到 engine 映射的原因
func (group *RouterGroup) addRoute(method string, comp string, handler HandlerFunc) {
	pattern := group.prefix + comp
	log.Printf("Route %4s - %s", method, pattern)
	group.engine.router.addRoute(method, pattern, handler)
}

// GET defines the method to add GET request
func (group *RouterGroup) GET(pattern string, handler HandlerFunc) {
	group.addRoute("GET", pattern, handler)
}

// POST defines the method to add POST request
func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
	group.addRoute("POST", pattern, handler)
}

//func (engine *Engine) addRoute(method string, pattern string, handler HandlerFunc) {
//	engine.router.addRoute(method, pattern, handler)
//}
//
//// GET defines the method to add GET request
//func (engine *Engine) GET(pattern string, handler HandlerFunc) {
//	engine.addRoute("GET", pattern, handler)
//}
//
//// POST defines the method to add POST request
//func (engine *Engine) POST(pattern string, handler HandlerFunc) {
//	engine.addRoute("POST", pattern, handler)
//}

// Run defines the method to start a http server
func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := newContext(w, req)
	engine.router.handle(c)
}
