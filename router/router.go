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
	radix      *radix.Radix
	notFound   types.Handler
	middleware []types.Middleware
	started    *atomic.Bool
	prefix     string
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
		started:  &atomic.Bool{},
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

	responder := h(req)
	responder.Respond(w, req)
}

func (r *Router) add(method string, handler types.Handler) {
	if r.started.Load() {
		panic(fmt.Sprintf("cannot register path: %s since the router is running", r.prefix))
	}

	// Apply route-specific middleware in reverse order at registration time
	h := handler
	for i := len(r.middleware) - 1; i >= 0; i-- {
		h = r.middleware[i](h)
	}

	if err := r.radix.AddRoute(method, r.prefix, h); err != nil {
		panic(fmt.Sprintf("%s %s: %v", method, r.prefix, err))
	}
}

// GET registers a handler for GET requests at the router's current prefix path.
// The prefix can include parameters (e.g., "/users/:id") and wildcards (e.g., "/files/*filepath").
// Panics if the route cannot be registered (e.g., conflicts with existing routes).
func (r *Router) GET(handler types.Handler) {
	r.add(http.MethodGet, handler)
}

// POST registers a handler for POST requests at the router's current prefix path.
// The prefix can include parameters (e.g., "/users/:id") and wildcards (e.g., "/files/*filepath").
// Panics if the route cannot be registered (e.g., conflicts with existing routes).
func (r *Router) POST(handler types.Handler) {
	r.add(http.MethodPost, handler)
}

// PUT registers a handler for PUT requests at the router's current prefix path.
// The prefix can include parameters (e.g., "/users/:id") and wildcards (e.g., "/files/*filepath").
// Panics if the route cannot be registered (e.g., conflicts with existing routes).
func (r *Router) PUT(handler types.Handler) {
	r.add(http.MethodPut, handler)
}

// DELETE registers a handler for DELETE requests at the router's current prefix path.
// The prefix can include parameters (e.g., "/users/:id") and wildcards (e.g., "/files/*filepath").
// Panics if the route cannot be registered (e.g., conflicts with existing routes).
func (r *Router) DELETE(handler types.Handler) {
	r.add(http.MethodDelete, handler)
}

// PATCH registers a handler for PATCH requests at the router's current prefix path.
// The prefix can include parameters (e.g., "/users/:id") and wildcards (e.g., "/files/*filepath").
// Panics if the route cannot be registered (e.g., conflicts with existing routes).
func (r *Router) PATCH(handler types.Handler) {
	r.add(http.MethodPatch, handler)
}

// HEAD registers a handler for HEAD requests at the router's current prefix path.
// The prefix can include parameters (e.g., "/users/:id") and wildcards (e.g., "/files/*filepath").
// Panics if the route cannot be registered (e.g., conflicts with existing routes).
func (r *Router) HEAD(handler types.Handler) {
	r.add(http.MethodHead, handler)
}

// OPTIONS registers a handler for OPTIONS requests at the router's current prefix path.
// The prefix can include parameters (e.g., "/users/:id") and wildcards (e.g., "/files/*filepath").
// Panics if the route cannot be registered (e.g., conflicts with existing routes).
func (r *Router) OPTIONS(handler types.Handler) {
	r.add(http.MethodOptions, handler)
}

// CONNECT registers a handler for CONNECT requests at the router's current prefix path.
// The prefix can include parameters (e.g., "/users/:id") and wildcards (e.g., "/files/*filepath").
// Panics if the route cannot be registered (e.g., conflicts with existing routes).
func (r *Router) CONNECT(handler types.Handler) {
	r.add(http.MethodConnect, handler)
}

// TRACE registers a handler for TRACE requests at the router's current prefix path.
// The prefix can include parameters (e.g., "/users/:id") and wildcards (e.g., "/files/*filepath").
// Panics if the route cannot be registered (e.g., conflicts with existing routes).
func (r *Router) TRACE(handler types.Handler) {
	r.add(http.MethodTrace, handler)
}

func (r *Router) shallowCopy() *Router {
	nr := Router{
		radix:      r.radix,
		notFound:   r.notFound,
		prefix:     r.prefix,
		started:    r.started,
		middleware: append([]types.Middleware{}, r.middleware...),
	}
	return &nr
}

// Use adds one or more middleware to the router's global middleware chain.
// Middleware is applied to all routes in the order it is registered.
// Multiple calls to Use will append middleware to the chain.
func (r *Router) Use(mws ...types.Middleware) *Router {
	nr := r.shallowCopy()
	nr.middleware = append(nr.middleware, mws...)
	return nr
}

func (r *Router) Prefix(segment string) *Router {
	if segment == "" {
		return r.shallowCopy() // no change
	}

	// trim trailing slash from existing prefix
	base := strings.TrimRight(r.prefix, "/")
	// trim leading slash from new segment
	seg := strings.TrimLeft(segment, "/")

	nr := r.shallowCopy()
	nr.prefix = base + "/" + seg
	return nr
}

// ServeStatic registers a handler to serve static files from the given filesystem.
// The router's current prefix determines the URL path where files will be served.
// For example, r.Prefix("/static").ServeStatic(os.DirFS("./static")) serves files from
// the ./static directory at /static/*.
// Automatically handles directory redirects and delegates to http.FileServer.
func (r *Router) ServeStatic(f fs.FS) {
	staticResponder := responders.NewStaticDirResponder(f, r.prefix)

	// Add wildcard pattern for file paths and register handler
	r.Prefix("/*fp").GET(func(req *http.Request) types.Responder {
		return staticResponder
	})
}
