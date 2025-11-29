package responders

import (
	"encoding/json"
	"net/http"
)

type jsonResponder struct {
	body   any
	status int
}

func JSONResponse(body any, status int) *jsonResponder {
	return &jsonResponder{body: body, status: status}
}

func (r *jsonResponder) Respond(w http.ResponseWriter, req *http.Request) {
	data, err := json.Marshal(r.body)
	if err != nil {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if r.status > 0 {
		w.WriteHeader(r.status)
	}
	w.Write(data)
}

type jsonErrorResponder struct {
	status int
	msg    string
}

func JSONErrorResponse(msg string, status int) *jsonErrorResponder {
	return &jsonErrorResponder{msg: msg, status: status}
}

type jsonError struct {
	Msg string `json:"msg"`
}

func (e *jsonErrorResponder) Respond(w http.ResponseWriter, req *http.Request) {
	data, err := json.Marshal(jsonError{Msg: e.msg})
	if err != nil {
		return
	}

	w.Header().Set("Content-Type", "application/problem+json")
	w.WriteHeader(e.status)
	w.Write(data)
}
