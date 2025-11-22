package radix

import (
	"fmt"
	"strings"

	"github.com/elmq0022/kami/types"
)

type Node struct {
	prefix       string
	children     []*Node
	paramName    string
	param        *Node
	wildcardName string
	wildcard     *Node
	terminal     map[string]types.Handler
}

type Radix struct {
	root *Node
}

func New() (*Radix, error) {
	r := Radix{root: &Node{}}

	// for _, route := range routes {
	// 	if len(route.Path) == 0 || route.Path[0] != '/' {
	// 		return nil, fmt.Errorf("path must start with '/'")
	// 	}

	// 	segments := pathSegments(route.Path)
	// 	if err := r.insert(route, r.root, segments, 0); err != nil {
	// 		return nil, err
	// 	}
	// }

	return &r, nil
}

func (r *Radix) AddRoute(method string, path string, handler types.Handler) error {
	if len(path) == 0 || path[0] != '/' {
		return fmt.Errorf("path must start with '/'")
	}
	route := types.Route{Method: method, Path: path, Handler: handler}
	segments := pathSegments(path)
	return r.insert(route, r.root, segments, 0)
}

func (r *Radix) insert(route types.Route, node *Node, segments []string, pos int) error {
	if pos >= len(segments) {
		if node.terminal == nil {
			node.terminal = make(map[string]types.Handler)
		}
		node.terminal[route.Method] = route.Handler
		return nil
	}

	seg := segments[pos]

	if len(seg) > 1 && seg[0] == ':' {
		if node.param == nil {
			node.param = &Node{paramName: seg[1:]}
			return r.insert(route, node.param, segments, pos+1)
		} else if node.param.paramName == seg[1:] {
			return r.insert(route, node.param, segments, pos+1)
		} else {
			return fmt.Errorf("parameter name conflict: existing '%s' vs new '%s' in path '%s'", node.param.paramName, seg[1:], route.Path)
		}
	}

	if len(seg) > 1 && seg[0] == '*' {
		if pos != len(segments)-1 {
			return fmt.Errorf("wildcard in non-terminal position in path '%s'", route.Path)
		}
		if node.wildcard == nil {
			node.wildcard = &Node{wildcardName: seg[1:]}
			return r.insert(route, node.wildcard, segments, pos+1)
		}
		return fmt.Errorf("multiple wildcards at same node for path '%s'", route.Path)
	}

	for _, child := range node.children {
		if child.prefix == seg {
			return r.insert(route, child, segments, pos+1)
		}
	}

	n := &Node{prefix: seg}
	node.children = append(node.children, n)
	return r.insert(route, n, segments, pos+1)
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

	if pos >= len(segments) {
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

	if node.wildcard != nil {
		params[node.wildcard.wildcardName] = strings.Join(segments[pos:], "/")
		h, ok := node.wildcard.terminal[method]
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
