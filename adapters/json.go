package adapters

import (
	"encoding/json"
	"net/http"

	"github.com/elmq0022/kami/types"
)

func JsonAdapter(w http.ResponseWriter, req *http.Request, handler types.Handler) {
	resp, err := handler(req)

	w.Header().Set("Content-Type", "application/json")

	if err != nil {
		// Handler returned an error; respond with 500 unless status provided
		if resp.Status == 0 {
			resp.Status = http.StatusInternalServerError
		}
		errorResponse := map[string]any{"error": err.Error()}
		data, jerr := json.Marshal(errorResponse)
		if jerr != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(resp.Status)
		_, _ = w.Write(data)
		return
	}

	if resp.Status == 0 {
		resp.Status = http.StatusOK
	}

	data, jerr := json.Marshal(resp.Body)
	if jerr != nil {
		http.Error(w, "failed to marshal response", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(resp.Status)
	_, _ = w.Write(data)
}
