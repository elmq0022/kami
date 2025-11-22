package router

import (
	"net/http"

	"github.com/elmq0022/kami/handlers"
	"github.com/elmq0022/kami/internal/radix"
	"github.com/elmq0022/kami/types"
)

type Router struct {
	adapter  types.Adapter
	radix    *radix.Radix
	notFound types.Handler
}

func New(adapter types.Adapter) (*Router, error) {
	rdx, err := radix.New()
	if err != nil {
		return nil, err
	}

	return &Router{
		adapter:  adapter,
		radix:    rdx,
		notFound: handlers.DefaultNotFoundHandler,
	}, nil
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	h, params, ok := r.radix.Lookup(req.Method, req.URL.Path)
	if !ok {
		h = r.notFound
		params = map[string]string{}
	}

	ctx := WithParams(req.Context(), params)
	req = req.WithContext(ctx)

	r.adapter(w, req, h)
}

func (r *Router) add(method, path string, handler types.Handler) error {
	return r.radix.AddRoute(method, path, handler)
}

func (r *Router) GET(path string, handler types.Handler) error {
	return r.add(http.MethodGet, path, handler)
}

func (r *Router) POST(path string, handler types.Handler) error {
	return r.add(http.MethodPost, path, handler)
}

func (r *Router) PUT(path string, handler types.Handler) error {
	return r.add(http.MethodPut, path, handler)
}

func (r *Router) DELETE(path string, handler types.Handler) error {
	return r.add(http.MethodDelete, path, handler)
}

func (r *Router) PATCH(path string, handler types.Handler) error {
	return r.add(http.MethodPatch, path, handler)
}
