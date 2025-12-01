//go:build examples

package main

import (
	"net/http"

	"github.com/elmq0022/kami/responders"
	"github.com/elmq0022/kami/router"
	"github.com/elmq0022/kami/types"
)

func main() {
	// initialize the  app
	r, err := router.New()
	if err != nil {
		panic(err)
	}

	// add routes
	r.GET("/", hello)
	r.GET("/user/:id", getUser)

	r.Run(":8080")
}

// handler must return a types.Responder
// here we use the JSONResponse constructor to create
// a new jsonResponder the body map[string]string is
// marshaled to an equivalet string and written to the body
func hello(r *http.Request) types.Responder {
	return responders.JSONResponse(
		map[string]string{
			"message": "Hello, World!",
		},
		http.StatusOK,
	)
}

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// same idea as above but here we marshal an annoted
// struct to be returned as json. this can be used
// to enforce more rigid types
func getUser(r *http.Request) types.Responder {
	params := router.GetParams(r.Context())
	id := params["id"]

	return responders.JSONResponse(
		User{
			ID:   id,
			Name: "John Doe",
		},
		http.StatusOK,
	)
}
