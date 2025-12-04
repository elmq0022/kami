//go:build examples

package main

import (
	"embed"
	"io/fs"

	"github.com/elmq0022/kami/router"
)

//go:embed web
var webFS embed.FS

func main() {
	r, err := router.New()
	if err != nil {
		panic(err)
	}

	// Strip the "web" directory prefix from the embedded FS
	web, err := fs.Sub(webFS, "web")
	if err != nil {
		panic(err)
	}

	// serve static files using the builder-style API
	// the prefix determines where files are served
	r.Prefix("/").ServeStatic(web)

	// run the app
	r.Run(":8080")
}
