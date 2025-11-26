package router

import (
	"fmt"
	"log"
	"net/http"

	"github.com/elmq0022/kami/handlers"
	"github.com/elmq0022/kami/internal/radix"
	"github.com/elmq0022/kami/types"
)

type Router struct {
	radix    *radix.Radix
	notFound types.Handler
	global   []types.Middleware
}

func New(opts ...Option) (*Router, error) {
	rdx, err := radix.New()
	if err != nil {
		return nil, err
	}

	r := &Router{
		radix:    rdx,
		notFound: handlers.DefaultNotFoundHandler,
	}

	for _, opt := range opts {
		opt(r)
	}

	return r, nil
}

func (r *Router) Run(port string) {
	log.Printf("Starting server on %s", port)
	if err := http.ListenAndServe(port, r); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("panic handling %s %s: %v", req.Method, req.URL.Path, err)
			http.Error(
				w,
				http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError,
			)
		}
	}()

	h, params, ok := r.radix.Lookup(req.Method, req.URL.Path)
	if !ok {
		h = r.notFound
		params = map[string]string{}
	}

	ctx := WithParams(req.Context(), params)
	req = req.WithContext(ctx)

	for i := len(r.global) - 1; i >= 0; i-- {
		h = r.global[i](h)
	}

	responder := h(req)
	responder.Respond(w)
}

func (r *Router) add(method, path string, handler types.Handler) {
	if err := r.radix.AddRoute(method, path, handler); err != nil {
		panic(fmt.Sprintf("%s %s: %v", method, path, err))
	}
}

func (r *Router) GET(path string, handler types.Handler) {
	r.add(http.MethodGet, path, handler)
}

func (r *Router) POST(path string, handler types.Handler) {
	r.add(http.MethodPost, path, handler)
}

func (r *Router) PUT(path string, handler types.Handler) {
	r.add(http.MethodPut, path, handler)
}

func (r *Router) DELETE(path string, handler types.Handler) {
	r.add(http.MethodDelete, path, handler)
}

func (r *Router) PATCH(path string, handler types.Handler) {
	r.add(http.MethodPatch, path, handler)
}

func (r *Router) HEAD(path string, handler types.Handler) {
	r.add(http.MethodHead, path, handler)
}

func (r *Router) OPTIONS(path string, handler types.Handler) {
	r.add(http.MethodOptions, path, handler)
}

func (r *Router) CONNECT(path string, handler types.Handler) {
	r.add(http.MethodConnect, path, handler)
}

func (r *Router) TRACE(path string, handler types.Handler) {
	r.add(http.MethodTrace, path, handler)
}

func (r *Router) Group(prefix string) SubRouter {
	return NewSubRouter(r, prefix)
}
