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
}

func (r *StaticDirectoryResponder) Respond(w http.ResponseWriter, req *http.Request) {
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

	http.ServeFile(w, req, fullPath)
}
