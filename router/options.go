package router

import "github.com/elmq0022/kami/types"

type Option func(r *Router)

func WithMiddleware() Option {
	return func(r *Router) {}
}

func WithNotFound(h types.Handler) Option {
	return func(r *Router) {
		r.notFound = h
	}
}

func WithLogger() Option {
	return func(r *Router) {}
}
