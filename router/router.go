// Package router provides HTTP routing functionality using a radix tree for efficient path matching.
// It supports path parameters, wildcards, middleware, and grouped routes.
package router

import (
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"strings"
	"sync/atomic"

	"github.com/elmq0022/kami/handlers"
	"github.com/elmq0022/kami/internal/radix"
	"github.com/elmq0022/kami/responders"
	"github.com/elmq0022/kami/types"
)

// Router is the main HTTP router that uses a radix tree for efficient route matching.
// It supports middleware, custom 404 handlers, and panic recovery.
type Router struct {
	radix    *radix.Radix
	notFound types.Handler
	global   []types.Middleware
	started  atomic.Bool
}

// New creates a new Router with the given options.
// Options can configure middleware, custom 404 handlers, and other router behavior.
// Returns an error if the underlying radix tree initialization fails.
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

// Run starts the HTTP server on the specified port.
// The port should be in the format ":8080" or "localhost:8080".
// This is a convenience method that calls http.ListenAndServe with the router as the handler.
// The function will block until the server fails to start or is shut down.
func (r *Router) Run(port string) {
	r.started.Store(true)
	log.Printf("Starting server on %s", port)
	if err := http.ListenAndServe(port, r); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

// ServeHTTP implements http.Handler, making Router compatible with the standard library.
// It performs route lookup, applies middleware, handles panics, and executes the matched handler.
// If no route matches, the configured notFound handler is used (defaults to a 404 response).
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.started.Store(true)

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
	responder.Respond(w, req)
}

func (r *Router) add(method, path string, handler types.Handler, middleware ...types.Middleware) {
	if r.started.Load() {
		panic(fmt.Sprintf("cannot register path: %s since the router is running", path))
	}

	// Apply route-specific middleware in reverse order at registration time
	h := handler
	for i := len(middleware) - 1; i >= 0; i-- {
		h = middleware[i](h)
	}

	if err := r.radix.AddRoute(method, path, h); err != nil {
		panic(fmt.Sprintf("%s %s: %v", method, path, err))
	}
}

// GET registers a handler for GET requests at the given path.
// Path can include parameters (e.g., "/users/:id") and wildcards (e.g., "/files/*filepath").
// Optional middleware can be provided as additional arguments and will be applied to this route only.
// Panics if the route cannot be registered (e.g., conflicts with existing routes).
func (r *Router) GET(path string, handler types.Handler, middleware ...types.Middleware) {
	r.add(http.MethodGet, path, handler, middleware...)
}

// POST registers a handler for POST requests at the given path.
// Path can include parameters (e.g., "/users/:id") and wildcards (e.g., "/files/*filepath").
// Optional middleware can be provided as additional arguments and will be applied to this route only.
// Panics if the route cannot be registered (e.g., conflicts with existing routes).
func (r *Router) POST(path string, handler types.Handler, middleware ...types.Middleware) {
	r.add(http.MethodPost, path, handler, middleware...)
}

// PUT registers a handler for PUT requests at the given path.
// Path can include parameters (e.g., "/users/:id") and wildcards (e.g., "/files/*filepath").
// Optional middleware can be provided as additional arguments and will be applied to this route only.
// Panics if the route cannot be registered (e.g., conflicts with existing routes).
func (r *Router) PUT(path string, handler types.Handler, middleware ...types.Middleware) {
	r.add(http.MethodPut, path, handler, middleware...)
}

// DELETE registers a handler for DELETE requests at the given path.
// Path can include parameters (e.g., "/users/:id") and wildcards (e.g., "/files/*filepath").
// Optional middleware can be provided as additional arguments and will be applied to this route only.
// Panics if the route cannot be registered (e.g., conflicts with existing routes).
func (r *Router) DELETE(path string, handler types.Handler, middleware ...types.Middleware) {
	r.add(http.MethodDelete, path, handler, middleware...)
}

// PATCH registers a handler for PATCH requests at the given path.
// Path can include parameters (e.g., "/users/:id") and wildcards (e.g., "/files/*filepath").
// Optional middleware can be provided as additional arguments and will be applied to this route only.
// Panics if the route cannot be registered (e.g., conflicts with existing routes).
func (r *Router) PATCH(path string, handler types.Handler, middleware ...types.Middleware) {
	r.add(http.MethodPatch, path, handler, middleware...)
}

// HEAD registers a handler for HEAD requests at the given path.
// Path can include parameters (e.g., "/users/:id") and wildcards (e.g., "/files/*filepath").
// Optional middleware can be provided as additional arguments and will be applied to this route only.
// Panics if the route cannot be registered (e.g., conflicts with existing routes).
func (r *Router) HEAD(path string, handler types.Handler, middleware ...types.Middleware) {
	r.add(http.MethodHead, path, handler, middleware...)
}

// OPTIONS registers a handler for OPTIONS requests at the given path.
// Path can include parameters (e.g., "/users/:id") and wildcards (e.g., "/files/*filepath").
// Optional middleware can be provided as additional arguments and will be applied to this route only.
// Panics if the route cannot be registered (e.g., conflicts with existing routes).
func (r *Router) OPTIONS(path string, handler types.Handler, middleware ...types.Middleware) {
	r.add(http.MethodOptions, path, handler, middleware...)
}

// CONNECT registers a handler for CONNECT requests at the given path.
// Path can include parameters (e.g., "/users/:id") and wildcards (e.g., "/files/*filepath").
// Optional middleware can be provided as additional arguments and will be applied to this route only.
// Panics if the route cannot be registered (e.g., conflicts with existing routes).
func (r *Router) CONNECT(path string, handler types.Handler, middleware ...types.Middleware) {
	r.add(http.MethodConnect, path, handler, middleware...)
}

// TRACE registers a handler for TRACE requests at the given path.
// Path can include parameters (e.g., "/users/:id") and wildcards (e.g., "/files/*filepath").
// Optional middleware can be provided as additional arguments and will be applied to this route only.
// Panics if the route cannot be registered (e.g., conflicts with existing routes).
func (r *Router) TRACE(path string, handler types.Handler, middleware ...types.Middleware) {
	r.add(http.MethodTrace, path, handler, middleware...)
}

// Group creates a SubRouter with the given path prefix.
// All routes registered on the SubRouter will be prefixed with this path.
// For example, Group("/api/v1") creates routes under /api/v1/*.
func (r *Router) Group(prefix string) SubRouter {
	return NewSubRouter(r, prefix)
}

// Use adds one or more middleware to the router's global middleware chain.
// Middleware is applied to all routes in the order it is registered.
// Multiple calls to Use will append middleware to the chain.
func (r *Router) Use(mw ...types.Middleware) {
	r.global = append(r.global, mw...)
}

// ServeStatic registers a handler to serve static files from the given filesystem.
// The prefix determines the URL path where files will be served.
// For example, ServeStatic(os.DirFS("./static"), "/static") serves files from
// the ./static directory at /static/*.
// Automatically handles directory redirects and delegates to http.FileServer.
func (r *Router) ServeStatic(f fs.FS, prefix string) {
	staticResponder := responders.NewStaticDirResponder(f, prefix)

	if !strings.HasSuffix(prefix, "/") {
		prefix += "/"
	}
	prefix += "*fp"

	// Wrap in closure if router expects a func
	r.GET(prefix, func(req *http.Request) types.Responder {
		return staticResponder
	})
}
