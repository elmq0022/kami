package radix

import (
	"fmt"
	"strings"

	"github.com/elmq0022/krillin/router"
)

type Node struct {
	prefix   string
	children []*Node
	terminal map[string]router.Handler
}

type Radix struct {
	root *Node
}

func New(routes []router.Route) (*Radix, error) {
	r := Radix{root: &Node{}}

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

func (r *Radix) addRoute(route router.Route, node *Node, segments []string, pos int) {
	if pos >= len(segments) {
		if node.terminal == nil {
			node.terminal = make(map[string]router.Handler)
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

	n := &Node{prefix: seg}
	node.children = append(node.children, n)
	r.addRoute(route, n, segments, pos+1)
}

func compress(node *Node) {
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

func (r *Radix) Lookup(method, path string) (router.Handler, bool) {
	root := r.root
	return lookup(root, method, path)
}

func lookup(node *Node, method, path string) (router.Handler, bool) {
	var zero router.Handler

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
		// check if the prefix matches and then ensure there is a complete match or a full segment is matched
		if strings.HasPrefix(path, child.prefix) && (len(path) == len(child.prefix) || path[len(child.prefix)] == '/') {
			h, ok := lookup(child, method, path[len(child.prefix):])
			return h, ok
		}
	}

	return zero, false
}
