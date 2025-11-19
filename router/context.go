package router

import "context"

type contextKey string

const paramsKey contextKey = "paramsKey"

func WithParams(ctx context.Context, params map[string]string) context.Context {
	return context.WithValue(ctx, paramsKey, params)
}

func GetParams(ctx context.Context) map[string]string {
	if p, ok := ctx.Value(paramsKey).(map[string]string); ok {
		return p
	}
	return map[string]string{}
}
