package handlers

import (
	"net/http"

	"github.com/elmq0022/kami/types"
)

func DefaultNotFoundHandler(r *http.Request) (types.Response, error) {
	return types.Response{Status: http.StatusNotFound, Body: "Not Found"}, nil
}
