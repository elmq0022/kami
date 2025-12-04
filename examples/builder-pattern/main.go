//go:build examples

package main

import (
	"fmt"
	"net/http"

	"github.com/elmq0022/kami/responders"
	"github.com/elmq0022/kami/router"
	"github.com/elmq0022/kami/types"
)

func main() {
	// Initialize the router
	r, err := router.New()
	if err != nil {
		panic(err)
	}

	// Public routes (no auth)
	r.Prefix("/").GET(home)
	r.Prefix("/health").GET(health)

	// API routes with authentication
	api := r.Prefix("/api").Use(authMiddleware)

	// API v1 routes
	v1 := api.Prefix("/v1")
	v1.Prefix("/status").GET(apiStatus)

	// User routes with logging middleware
	users := v1.Prefix("/users").Use(loggingMiddleware)
	users.GET(listUsers)
	users.Prefix("/:id").GET(getUser)

	// Admin routes with additional admin middleware
	admin := v1.Prefix("/admin").Use(adminMiddleware)
	admin.Prefix("/users").GET(listAllUsers)
	admin.Prefix("/users/:id").DELETE(deleteUser)

	fmt.Println("Server starting on :8080")
	fmt.Println("Try these endpoints:")
	fmt.Println("  GET  /")
	fmt.Println("  GET  /health")
	fmt.Println("  GET  /api/v1/status")
	fmt.Println("  GET  /api/v1/users")
	fmt.Println("  GET  /api/v1/users/:id")
	fmt.Println("  GET  /api/v1/admin/users")
	fmt.Println("  DELETE /api/v1/admin/users/:id")

	r.Run(":8080")
}

// Middleware

func authMiddleware(next types.Handler) types.Handler {
	return func(r *http.Request) types.Responder {
		// In a real app, you'd check headers, tokens, etc.
		fmt.Println("[Auth] Checking authentication...")
		return next(r)
	}
}

func loggingMiddleware(next types.Handler) types.Handler {
	return func(r *http.Request) types.Responder {
		fmt.Printf("[Log] %s %s\n", r.Method, r.URL.Path)
		return next(r)
	}
}

func adminMiddleware(next types.Handler) types.Handler {
	return func(r *http.Request) types.Responder {
		// In a real app, you'd check admin privileges
		fmt.Println("[Admin] Checking admin privileges...")
		return next(r)
	}
}

// Handlers

func home(r *http.Request) types.Responder {
	return responders.JSONResponse(
		map[string]string{
			"message": "Welcome to the API!",
			"version": "1.0.0",
		},
		http.StatusOK,
	)
}

func health(r *http.Request) types.Responder {
	return responders.JSONResponse(
		map[string]string{
			"status": "healthy",
		},
		http.StatusOK,
	)
}

func apiStatus(r *http.Request) types.Responder {
	return responders.JSONResponse(
		map[string]string{
			"api":     "v1",
			"status":  "operational",
			"message": "This endpoint has auth middleware",
		},
		http.StatusOK,
	)
}

func listUsers(r *http.Request) types.Responder {
	return responders.JSONResponse(
		map[string]interface{}{
			"users": []map[string]string{
				{"id": "1", "name": "Alice"},
				{"id": "2", "name": "Bob"},
			},
			"message": "This endpoint has auth + logging middleware",
		},
		http.StatusOK,
	)
}

func getUser(r *http.Request) types.Responder {
	params := router.GetParams(r.Context())
	id := params["id"]

	return responders.JSONResponse(
		map[string]interface{}{
			"user": map[string]string{
				"id":   id,
				"name": "User " + id,
			},
			"message": "This endpoint has auth + logging middleware",
		},
		http.StatusOK,
	)
}

func listAllUsers(r *http.Request) types.Responder {
	return responders.JSONResponse(
		map[string]interface{}{
			"users": []map[string]string{
				{"id": "1", "name": "Alice", "role": "user"},
				{"id": "2", "name": "Bob", "role": "admin"},
				{"id": "3", "name": "Charlie", "role": "user"},
			},
			"message": "This endpoint has auth + admin middleware",
		},
		http.StatusOK,
	)
}

func deleteUser(r *http.Request) types.Responder {
	params := router.GetParams(r.Context())
	id := params["id"]

	return responders.JSONResponse(
		map[string]interface{}{
			"deleted": id,
			"message": "User deleted (this endpoint has auth + admin middleware)",
		},
		http.StatusOK,
	)
}
