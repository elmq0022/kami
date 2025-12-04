# KAMI

## Description
A small but capable micro web framework for the Go programming language.

## Objectives
The library is primarily aimed at microservices that back frontend applications consuming JSON via JavaScript's Fetch API.
The framework also supports serving static files from directories.
Web templates are also planned.

## Philosophy
The author aims to keep the library small enough that reading the code and a few examples can serve as the documentation.

## Getting Started

### Installing Task

This project uses [Task](https://taskfile.dev/) as a task runner. To install it:

**macOS/Linux:**
```bash
sh -c "$(curl --location https://taskfile.dev/install.sh)" -- -d -b ~/.local/bin
```

**macOS (Homebrew):**
```bash
brew install go-task/tap/go-task
```

**Other installation methods:**
See the [official Task installation guide](https://taskfile.dev/installation/) for additional options including package managers for various platforms.

## Examples

The project includes example applications demonstrating various features. See [examples/README.md](examples/README.md) for detailed documentation.

### Quickstart Example
A basic example showing routing and JSON responses.

**Run with:**
```bash
task run:example:quickstart
```

**Features demonstrated:**
- Builder-style API with `Prefix()` chaining
- Basic routing with `GET` requests
- URL parameters (`:id`)
- JSON responses using `responders.JSONResponse`
- Handler functions returning `types.Responder`

**Source:** [examples/quickstart/main.go](examples/quickstart/main.go)

### Static File Serving Example
Demonstrates serving static files from an embedded filesystem.

**Run with:**
```bash
task run:example:static
```

**Features demonstrated:**
- Serving static files with `ServeStatic`
- Using Go's `embed.FS` to bundle assets
- Builder-style API for route registration

**Source:** [examples/static/main.go](examples/static/main.go)

### Builder Pattern Example
Advanced example showcasing middleware composition and route organization.

**Run with:**
```bash
task run:example:builder
```

**Features demonstrated:**
- Nested route groups with `Prefix()` chaining
- Middleware composition with `Use()`
- Multiple middleware layers (auth, logging, admin)
- Public vs protected routes
- Immutable router pattern

**Source:** [examples/builder-pattern/main.go](examples/builder-pattern/main.go)

## Quick Start

Here's a minimal example to get you started:

```go
package main

import (
    "net/http"

    "github.com/elmq0022/kami/responders"
    "github.com/elmq0022/kami/router"
    "github.com/elmq0022/kami/types"
)

func main() {
    // Create a new router
    r, _ := router.New()

    // Register routes using the builder-style API
    r.Prefix("/").GET(hello)
    r.Prefix("/user/:id").GET(getUser)

    // Start the server
    r.Run(":8080")
}

func hello(r *http.Request) types.Responder {
    return responders.JSONResponse(
        map[string]string{"message": "Hello, World!"},
        http.StatusOK,
    )
}

func getUser(r *http.Request) types.Responder {
    params := router.GetParams(r.Context())
    id := params["id"]

    return responders.JSONResponse(
        map[string]string{"id": id, "name": "John Doe"},
        http.StatusOK,
    )
}
```

## Usage

### Builder-Style API

Kami uses a builder-style API where `Prefix()` and `Use()` return new router instances. This allows for clean route organization and middleware composition:

```go
r, _ := router.New()

// Chain prefixes to build nested routes
r.Prefix("/api").Prefix("/v1").Prefix("/users").GET(listUsersHandler)

// Combine prefixes and middleware
api := r.Prefix("/api").Use(authMiddleware)
api.Prefix("/users").GET(listUsersHandler)

// Routes are registered with the accumulated prefix
users := r.Prefix("/api").Prefix("/users")
users.GET(listUsersHandler)           // GET /api/users
users.Prefix("/:id").GET(getUserHandler)  // GET /api/users/:id
```

### Routing Paths

- Parameters are defined with a leading colon `:`
    - The router disallows path prefixes followed by a different parameter name. For example, registering both of these paths would lead to an error:
    ```
    /foo/bar/:buzz
    /foo/bar/:bazz
    ```

- Wildcards are defined with a leading asterisk `*`

- The match precedence for a path is:
  `static` → `:parameter` → `*wildcard`

### Context Parameters

- Any values read from the URL are stored in the request context
- A `map[string]string` of parameter value key-value pairs can be retrieved with `GetParams(req.Context())`
- If there are no params, expect an empty `map[string]string`
- Users should check that a value exists in the map using the standard Go idiom: `val, exists := params[key]`


### Middleware

The framework uses a builder-style API for composing middleware. Middleware is accumulated as you chain routers together and applied at **registration time**.

#### Adding Middleware

Middleware is added using the `Use()` method, which returns a new router instance with the middleware added:

```go
r, _ := router.New()

// Add middleware - returns a new router with middleware
r = r.Use(router.Logger)
r = r.Use(myCustomMiddleware1, myCustomMiddleware2)
```

**Important:** `Use()` returns a **new router instance**. You must assign the result if you want to use that router.

#### Middleware Composition

Middleware accumulates as you chain routers. Child routers inherit their parent's middleware:

```go
r, _ := router.New()

// Create a router with auth middleware
api := r.Prefix("/api").Use(authMiddleware)

// Create a child router with additional logging
users := api.Prefix("/users").Use(loggingMiddleware)

// Registers at /api/users with both auth and logging middleware
users.GET(listUsersHandler)
```

#### Route Groups

Use the builder pattern to organize routes into groups with different middleware:

```go
r, _ := router.New()

// Public routes (no middleware)
r.Prefix("/").GET(homeHandler)
r.Prefix("/health").GET(healthHandler)

// Protected API routes
api := r.Prefix("/api").Use(authMiddleware)
api.Prefix("/users").GET(listUsersHandler)
api.Prefix("/users/:id").GET(getUserHandler)

// Admin routes with additional middleware
admin := api.Prefix("/admin").Use(adminMiddleware)
admin.Prefix("/users").GET(listAllUsersHandler)
```

#### Execution Order

Middleware executes in the order it was added via `Use()`:

```go
r := r.Use(logger)           // 1st
r = r.Use(cors)              // 2nd
r = r.Use(auth)              // 3rd

// Execution order: logger -> cors -> auth -> handler
```

The middleware signature is:

```go
func(next types.Handler) types.Handler
```

#### Built-in Middleware

- `router.Logger` - Logs each request with method, path, status code, and duration

#### Key Principles

1. **Immutability**: `Use()` and `Prefix()` return new router instances
2. **Inheritance**: Child routers inherit parent middleware
3. **Isolation**: Sibling routers don't affect each other
4. **Registration-time composition**: Middleware is applied when routes are registered

### Path Registration

The addition of a path mutates the radix tree used for lookups and is NOT thread-safe.
The expectation is that routes will be registered prior to the server performing any path lookups.
Lookups are read-only and therefore thread-safe.

