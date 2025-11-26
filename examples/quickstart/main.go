package main

import (
	"net/http"

	"github.com/elmq0022/kami/responders.go"
	"github.com/elmq0022/kami/router"
	"github.com/elmq0022/kami/types"
)

func main() {
	r, err := router.New()
	if err != nil {
		panic(err)
	}

	r.GET("/", hello)
	r.GET("/user/:id", getUser)

	http.ListenAndServe(":8080", r)
}

func hello(r *http.Request) types.Responder {
	return &responders.JSONResponder{
		Body: map[string]string{
			"message": "Hello, World!",
		},
		Status: http.StatusOK,
	}
}

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func getUser(r *http.Request) types.Responder {
	params := router.GetParams(r.Context())
	id := params["id"]

	return &responders.JSONResponder{
		Body: User{
			ID:   id,
			Name: "John Doe",
		},
		Status: http.StatusOK,
	}
}
