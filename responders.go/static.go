package responders

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type StaticDirectoryResponder struct {
	BaseDir  string
	FilePath string
	Request  *http.Request
}

func (r *StaticDirectoryResponder) Respond(w http.ResponseWriter) {
	cleanPath := filepath.Clean("/" + r.FilePath)[1:]

	if strings.Contains(cleanPath, "..") {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	fullPath := filepath.Join(r.BaseDir, cleanPath)

	stat, err := os.Stat(fullPath)
	if err == nil && stat.IsDir() {
		fullPath = filepath.Join(fullPath, "index.html")
	}

	http.ServeFile(w, r.Request, fullPath)
}
