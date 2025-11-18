package radix

import (
	"fmt"
	"strings"

	"github.com/elmq0022/krillin/types"
)

type Node struct {
	prefix    string
	children  []*Node
	paramName string
	param     *Node
	terminal  map[string]types.Handler
}

type Radix struct {
	root *Node
}

func New(routes types.Routes) (*Radix, error) {
	r := Radix{root: &Node{}}

	for _, route := range routes {
		if len(route.Path) == 0 || route.Path[0] != '/' {
			return nil, fmt.Errorf("path must start with '/'")
		}

		segments := pathSegments(route.Path)
		r.addRoute(route, r.root, segments, 0)
	}

	return &r, nil
}

func (r *Radix) addRoute(route types.Route, node *Node, segments []string, pos int) {
	if pos >= len(segments) {
		if node.terminal == nil {
			node.terminal = make(map[string]types.Handler)
		}
		node.terminal[route.Method] = route.Handler
		return
	}

	seg := segments[pos]

	if len(seg) > 0 && seg[0] == ':' {
		if node.param == nil {
			node.param = &Node{paramName: seg[1:]}
			r.addRoute(route, node.param, segments, pos+1)
			return
		} else if node.param.paramName == seg[1:] {
			r.addRoute(route, node.param, segments, pos+1)
			return
		} else {
			panic("bad")
		}
	}

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

func (r *Radix) Lookup(method, path string) (types.Handler, map[string]string, bool) {
	root := r.root
	segments := pathSegments(path)
	params := make(map[string]string)
	handler, ok := lookup(root, method, segments, 0, params)
	return handler, params, ok
}

func lookup(node *Node, method string, segments []string, pos int, params map[string]string) (types.Handler, bool) {
	var zero types.Handler

	if node == nil {
		return zero, false
	}

	if pos >= len(segments) || len(segments) == 0 {
		handler, ok := node.terminal[method]
		return handler, ok
	}

	for _, child := range node.children {
		if segments[pos] == child.prefix {
			h, ok := lookup(child, method, segments, pos+1, params)
			return h, ok
		}
	}

	if node.param != nil {
		params[node.param.paramName] = segments[pos]
		h, ok := lookup(node.param, method, segments, pos+1, params)
		return h, ok
	}

	return zero, false
}

func pathSegments(path string) []string {
	segments := strings.Split(path, "/")

	p := 0
	for _, segment := range segments {
		if segment != "" {
			segments[p] = segment
			p++
		}
	}
	return segments[:p]
}
