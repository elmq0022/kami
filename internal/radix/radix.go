package radix

import (
	"fmt"
	"strings"

	"github.com/elmq0022/krillin/router"
)

type Node[T any] struct {
	prefix   string
	children []*Node[T]
	terminal map[string]T
}

type Radix[T any] struct {
	root *Node[T]
}

func New[T any](routes []router.Route[T]) (*Radix[T], error) {
	r := Radix[T]{root: &Node[T]{}}

	for _, route := range routes {
		if len(route.Path) == 0 || route.Path[0] != '/' {
			return nil, fmt.Errorf("path must start with '/'")
		}

		segments := strings.Split(route.Path, "/")[1:]
		r.addRoute(route, r.root, segments, 0)
	}

	compress(r.root)
	return &r, nil
}

func (r *Radix[T]) addRoute(route router.Route[T], node *Node[T], segments []string, pos int) {
	if pos >= len(segments) {
		if node.terminal == nil {
			node.terminal = make(map[string]T)
		}
		node.terminal[route.Method] = route.Handler
		return
	}

	seg := segments[pos]

	for _, child := range node.children {
		if child.prefix == seg {
			r.addRoute(route, child, segments, pos+1)
			return
		}
	}

	n := &Node[T]{prefix: seg}
	node.children = append(node.children, n)
	r.addRoute(route, n, segments, pos+1)
}

func compress[T any](node *Node[T]) {
	for i := range node.children {
		compress(node.children[i])
	}

	if node.prefix == "" {
		return
	}

	if len(node.children) == 1 && node.terminal == nil {
		child := node.children[0]
		node.prefix = node.prefix + "/" + child.prefix
		node.terminal = child.terminal
		node.children = child.children
	}
}

func (r *Radix[T]) Lookup(method, path string) (T, bool) {
	root := r.root
	return lookup(root, method, path)
}

func lookup[T any](node *Node[T], method, path string) (T, bool) {
	var zero T

	if node == nil {
		return zero, false
	}

	if len(path) > 0 && path[0] == '/' {
		path = path[1:]
	}

	if path == "" {
		handler, ok := node.terminal[method]
		return handler, ok
	}

	for _, child := range node.children {
		if strings.HasPrefix(path, child.prefix) {
			path := path[len(child.prefix):]
			h, ok := lookup(child, method, path)
			return h, ok
		}
	}

	return zero, false
}

func (r *Radix[T]) ChildNPrefix(n int) string {
	return "/" + r.root.children[n].prefix
}
