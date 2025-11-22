package types

import (
	"net/http"
)

type Response struct {
	Status int
	Body   any
}
type Handler func(req *http.Request) (Response, error)
type Routes []Route
type Adapter func(http.ResponseWriter, *http.Request, Handler)

type Route struct {
	Method  string
	Path    string
	Handler Handler
}
