package handlers

import (
	"net/http"

	"github.com/elmq0022/kami/types"
)

type DefaultNotFoundResponder struct {
	Status int
	Body   string
}

func (dnf *DefaultNotFoundResponder) Respond(w http.ResponseWriter) {
	w.WriteHeader(dnf.Status)
	w.Write([]byte(dnf.Body))
}

func DefaultNotFoundHandler(r *http.Request) types.Responder {
	return &DefaultNotFoundResponder{Status: http.StatusNotFound, Body: "Not Found"}
}
