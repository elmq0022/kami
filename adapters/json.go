package adapters

import (
	"encoding/json"
	"net/http"

	"github.com/elmq0022/kami/types"
)

func JsonAdapter(w http.ResponseWriter, req *http.Request, handler types.Handler) {
	status, result, _ := handler(req)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	data, _ := json.Marshal(result)
	w.Write([]byte(data))
}
