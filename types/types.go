package types

import (
	"net/http"
)

type Renderable interface {
	Render(w http.ResponseWriter)
}

type Response struct {
	Status int
	Body   any
}

type Middleware func(h Handler) Handler
type Handler func(req *http.Request) Renderable
type Routes []Route

type Route struct {
	Method  string
	Path    string
	Handler Handler
}
