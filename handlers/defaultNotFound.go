package handlers

import (
	"net/http"

	"github.com/elmq0022/kami/types"
)

type DefaultNotFoundRenerable struct {
	Status int
	Body   string
}

func (dnf *DefaultNotFoundRenerable) Render(w http.ResponseWriter) {
	w.WriteHeader(dnf.Status)
	w.Write([]byte(dnf.Body))
}

func DefaultNotFoundHandler(r *http.Request) types.Renderable {
	return &DefaultNotFoundRenerable{Status: http.StatusNotFound, Body: "Not Found"}
}
