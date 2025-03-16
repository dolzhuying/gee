package gee

import "net/http"

type HandlerFunc func(*Context)

type(
	RouterGroup struct{
		prefix string
		middlewares []HandlerFunc
		parent *RouterGroup
		engine *Engine
	}

	Engine struct{
		*RouterGroup
		router *router
		groups []*RouterGroup
	}
)

func New() *Engine{
	engine:=&Engine{router: newRouter()}
	engine.RouterGroup=&RouterGroup{engine: engine}
	engine.groups=[]*RouterGroup{engine.RouterGroup}
	return engine
}

func (group*RouterGroup) Group(prefix string)*RouterGroup{
	engine:=group.engine
	newGroup:=&RouterGroup{
		prefix:group.prefix+prefix,
		parent:group,
		engine:engine,
	}
	engine.groups=append(engine.groups,newGroup)
	return newGroup
} 

func (group*RouterGroup) addRoute(method,comp string,handler HandlerFunc){
	pattern:=group.prefix+comp
	group.engine.router.addRoute(method,pattern,handler)
}

func (group *RouterGroup) GET(pattern string, handler HandlerFunc) {
	group.addRoute("GET", pattern, handler)
}

func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
	group.addRoute("POST", pattern, handler)
}

func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := newContext(w, req)
	engine.router.handle(c)
}
