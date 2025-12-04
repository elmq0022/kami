# Kami Examples

This directory contains example applications demonstrating the Kami HTTP router framework.

## Building Examples

All examples use the `examples` build tag:

```bash
go build -tags examples ./examples/quickstart
go build -tags examples ./examples/static
go build -tags examples ./examples/builder-pattern
```

## Examples

### 1. Quickstart (`quickstart/`)

Basic example showing simple route registration:

```go
r.Prefix("/").GET(hello)
r.Prefix("/user/:id").GET(getUser)
```

**Run it:**
```bash
go run -tags examples ./examples/quickstart
curl http://localhost:8080/
curl http://localhost:8080/user/123
```

### 2. Static Files (`static/`)

Example of serving static files from an embedded filesystem:

```go
r.Prefix("/").ServeStatic(web)
```

**Run it:**
```bash
go run -tags examples ./examples/static
curl http://localhost:8080/
```

### 3. Builder Pattern (`builder-pattern/`)

Advanced example showcasing the builder-style API with:
- Multiple route groups
- Nested prefixes
- Middleware composition
- Authentication and authorization layers

```go
// Create route groups with middleware
api := r.Prefix("/api").Use(authMiddleware)
v1 := api.Prefix("/v1")

// Add more middleware to specific groups
users := v1.Prefix("/users").Use(loggingMiddleware)
users.GET(listUsers)

// Create admin routes with additional middleware
admin := v1.Prefix("/admin").Use(adminMiddleware)
admin.Prefix("/users").GET(listAllUsers)
```

**Run it:**
```bash
go run -tags examples ./examples/builder-pattern
curl http://localhost:8080/api/v1/users
curl http://localhost:8080/api/v1/admin/users
```

## Key Concepts

### Builder-Style API

The router uses a builder pattern where methods return new router instances:

```go
// Each call returns a new router
api := r.Prefix("/api")
v1 := api.Prefix("/v1")
users := v1.Prefix("/users")

// Routers are immutable - creating children doesn't modify parents
parent := r.Use(middleware1)
child := parent.Use(middleware2)  // parent still only has middleware1
```

### Prefix Chaining

Prefixes can be chained to build nested route structures:

```go
r.Prefix("/api").Prefix("/v1").Prefix("/users/:id").GET(handler)
// Registers route at: /api/v1/users/:id
```

### Middleware Composition

Middleware accumulates as you chain routers:

```go
r1 := r.Use(authMiddleware)              // Has: auth
r2 := r1.Use(loggingMiddleware)          // Has: auth + logging
r3 := r2.Use(validationMiddleware)       // Has: auth + logging + validation
```

### Route Registration

Routes are registered by calling HTTP method functions on a prefixed router:

```go
// The prefix determines the route path
r.Prefix("/users").GET(listUsers)        // GET /users
r.Prefix("/users/:id").GET(getUser)      // GET /users/:id
r.Prefix("/users/:id").DELETE(deleteUser) // DELETE /users/:id
```

## Advanced Patterns

### Organizing by Feature

```go
// User routes
userAPI := r.Prefix("/api/users").Use(authMiddleware)
userAPI.GET(listUsers)
userAPI.Prefix("/:id").GET(getUser)
userAPI.POST(createUser)

// Product routes
productAPI := r.Prefix("/api/products").Use(authMiddleware)
productAPI.GET(listProducts)
productAPI.Prefix("/:id").GET(getProduct)
```

### Versioned APIs

```go
v1 := r.Prefix("/api/v1")
v1.Prefix("/users").GET(listUsersV1)

v2 := r.Prefix("/api/v2")
v2.Prefix("/users").GET(listUsersV2)
```

### Mixing Public and Protected Routes

```go
// Public routes
r.Prefix("/").GET(home)
r.Prefix("/health").GET(health)

// Protected routes with auth
protected := r.Use(authMiddleware)
protected.Prefix("/api/users").GET(listUsers)
protected.Prefix("/api/profile").GET(getProfile)
```
