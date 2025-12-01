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

	// serve static files by passing an fs.FS base directory
	// and the route base rout to r.ServeStatic
	r.ServeStatic(web, "/web/")

	// run the app
	r.Run(":8080")
}
