package responders

import (
	"io/fs"
	"net/http"
	"strings"
)

type staticDirectoryResponder struct {
	handler http.Handler
}

func NewStaticDirResponder(f fs.FS, prefix string) *staticDirectoryResponder {

	if prefix != "" && !strings.HasSuffix(prefix, "/") {
		prefix += "/"
	}

	fsHandler := http.FileServer(http.FS(f))
	if prefix != "" {
		fsHandler = http.StripPrefix(prefix, fsHandler)
	}

	return &staticDirectoryResponder{
		handler: fsHandler,
	}
}

func (r *staticDirectoryResponder) Respond(w http.ResponseWriter, req *http.Request) {
	r.handler.ServeHTTP(w, req)
}
