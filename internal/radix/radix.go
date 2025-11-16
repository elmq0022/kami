package radix

import (
	"fmt"

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
		rs := []rune(route.Path)
		r.addRoute(route, r.root, rs, 0)
	}

	compress(r.root)

	return &r, nil
}

func (r *Radix[T]) addRoute(route router.Route[T], root *Node[T], rs []rune, pos int) {
	if pos >= len(rs) {
		if root.terminal == nil {
			root.terminal = make(map[string]T)
		}
		root.terminal[route.Method] = route.Handler
		return
	}

	c := rs[pos]
	for _, node := range root.children {
		if node.prefix == string(c) {
			r.addRoute(route, node, rs, pos+1)
			return
		}
	}

	n := &Node[T]{prefix: string(c)}
	root.children = append(root.children, n)
	r.addRoute(route, n, rs, pos+1)
}

func compress[T any](n *Node[T]) {

	for i := range n.children {
		compress(n.children[i])
	}

	if len(n.children) == 1 && n.terminal == nil {
		child := n.children[0]

		n.prefix += child.prefix
		n.terminal = child.terminal
		n.children = child.children
	}
}

func (r *Radix[T]) Lookup(method, path string) (T, bool) {
	return *new(T), true
}

func (r *Radix[T]) FirstPrefix() string {
	return r.root.prefix
}
