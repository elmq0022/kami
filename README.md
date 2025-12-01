# KAMI

## Description
A small but capable micro web framework for the Go programming language.

## Objectives
The library is primarily aimed at microservices that back frontend applications consuming JSON via JavaScript's Fetch API.
The framework also supports serving static files from directories.
Web templates are also planned.

## Philosophy
The author aims to keep the library small enough that reading the code and a few examples can serve as the documentation.

## Usage

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

Middleware can be added to the router using the `Use` method, which accepts one or more middleware functions:

```go
r := router.New()
r.Use(router.Logger)
r.Use(myCustomMiddleware1, myCustomMiddleware2)
```

Middleware is applied to all routes in the order it is registered. The middleware signature is:

```go
func(next types.Handler) types.Handler
```

Built-in middleware:
- `router.Logger` - Logs each request with method, path, status code, and duration

### Path Registration

The addition of a path mutates the radix tree used for lookups and is NOT thread-safe.
The expectation is that routes will be registered prior to the server performing any path lookups.
Lookups are read-only and therefore thread-safe.

