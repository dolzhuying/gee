package gee

import (
	"fmt"
	"net/http"
)

type HandlerFunc func(http.ResponseWriter,*http.Request)

type Engine struct{
	router map[string]HandlerFunc
}

func New() *Engine{
	return &Engine{router:make(map[string]HandlerFunc)}
}

func (e *Engine) addroute(method string,pattern string,handler HandlerFunc){
	key:=method+"-"+pattern
	e.router[key]=handler
}

func (e *Engine) GET(pattern string,handler HandlerFunc){
	e.addroute("GET",pattern,handler)
}

func (e *Engine) POST(pattern string,handler HandlerFunc){
	e.addroute("POST",pattern,handler)
} 

func (e*Engine) Run(addr string) (err error){
	return http.ListenAndServe(addr,e)
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	key := r.Method + "-" + r.URL.Path
	if handler, ok := e.router[key]; ok {
		handler(w, r)
	} else {
		fmt.Fprintf(w, "404 NOT FOUND: %s\n", r.URL)
	}
}