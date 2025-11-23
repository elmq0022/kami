package router

import "github.com/elmq0022/kami/types"

type Option func(r *Router)

func WithMiddleware(mw types.Middleware) Option {
	return func(r *Router) {
		r.global = append(r.global, mw)
	}
}

func WithNotFound(h types.Handler) Option {
	return func(r *Router) {
		r.notFound = h
	}
}

func WithLogger() Option {
	return func(r *Router) {
		// TODO: implement this.
	}
}
