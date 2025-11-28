package responders

import (
	"encoding/json"
	"net/http"
)

type JSONResponder struct {
	Body   any
	Status int
}

func (r *JSONResponder) Respond(w http.ResponseWriter, req *http.Request) {
	data, err := json.Marshal(r.Body)
	if err != nil {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if r.Status > 0 {
		w.WriteHeader(r.Status)
	}
	w.Write(data)
}

type JSONErrorResponder struct {
	Status int
	Msg    string
}

type JSONError struct {
	Msg string `json:"msg"`
}

func (e *JSONErrorResponder) Respond(w http.ResponseWriter, req *http.Request) {
	data, err := json.Marshal(JSONError{Msg: e.Msg})
	if err != nil {
		return
	}

	w.Header().Set("Content-Type", "application/problem+json")
	w.WriteHeader(e.Status)
	w.Write(data)
}
