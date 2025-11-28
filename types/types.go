package types

import (
	"net/http"
)

type Responder interface {
	Respond(w http.ResponseWriter, r *http.Request)
}

type Response struct {
	Status int
	Body   any
}

type Middleware func(h Handler) Handler
type Handler func(req *http.Request) Responder
type Routes []Route

type Route struct {
	Method  string
	Path    string
	Handler Handler
}
