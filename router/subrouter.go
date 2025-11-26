package router

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/elmq0022/kami/types"
)

type SubRouter struct {
	r      *Router
	prefix string
}

func NewSubRouter(r *Router, prefix string) SubRouter {
	if len(prefix) < 1 {
		panic("prefix cannot be an empty string")
	}

	if prefix[0] != '/' {
		panic(fmt.Sprintf("prefix %s does not start with a '/'", prefix))
	}

	trimedPrefix := strings.TrimRight(prefix, "/")
	if trimedPrefix[len(trimedPrefix)-1] == '*' {
		panic(fmt.Sprintf("prefix %s cannot end in a wildcard '*'", prefix))
	}

	return SubRouter{
		r:      r,
		prefix: trimedPrefix,
	}
}

func (s *SubRouter) add(method, path string, handler types.Handler) {
	fullPath := s.prefix + "/" + strings.TrimLeft(path, "/")
	s.r.add(method, fullPath, handler)
}

func (s *SubRouter) GET(path string, handler types.Handler) {
	s.add(http.MethodGet, path, handler)
}

func (s *SubRouter) POST(path string, handler types.Handler) {
	s.add(http.MethodPost, path, handler)
}

func (s *SubRouter) PUT(path string, handler types.Handler) {
	s.add(http.MethodPut, path, handler)
}

func (s *SubRouter) PATCH(path string, handler types.Handler) {
	s.add(http.MethodPatch, path, handler)
}

func (s *SubRouter) DELETE(path string, handler types.Handler) {
	s.add(http.MethodDelete, path, handler)
}

func (s *SubRouter) HEAD(path string, handler types.Handler) {
	s.add(http.MethodHead, path, handler)
}

func (s *SubRouter) OPTIONS(path string, handler types.Handler) {
	s.add(http.MethodOptions, path, handler)
}
