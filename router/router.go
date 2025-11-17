package router

import (
	"net/http"
)

type Handler func(req *http.Request) (int, any, error)
type Routes []Route
type Adapter func(http.ResponseWriter, *http.Request, Handler)

type Route struct {
	Method  string
	Path    string
	Handler Handler
}

type Router struct {
	routes  []Route
	adapter Adapter
}

func New(routes []Route, processor Adapter) *Router {
	return &Router{
		routes:  routes,
		adapter: processor,
	}
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	for _, route := range r.routes {
		if route.Method == req.Method && route.Path == req.URL.Path {
			r.adapter(w, req, route.Handler)
			return
		}
	}
	http.NotFound(w, req)
}
