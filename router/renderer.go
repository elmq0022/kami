package router

import (
	"encoding/json"
	"net/http"
)

type JsonResponse struct {
	Status int
	Data   any
}

type HtmlTemplateResponse struct {
	Status   int
	Data     any
	Template string
}

type RawHtmlResponse struct {
	Status int
	Html   string
}

func DefaultRenderer(w http.ResponseWriter, r *http.Request, response any) {
	switch v := response.(type) {
	case JsonResponse:
		w.Header().Set("Content-Type", "application/json")
		data, err := json.Marshal(v.Data)
		if err != nil {
			http.Error(w, "failed to marshal response", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(v.Status)
		_, _ = w.Write(data)
	default:
	}
}
