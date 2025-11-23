# KAMI

## Description
A small but capable router library for the go programming language that speaks JSON.

## Objectives
The library is primarily aimed at microservices that back a frontend applications consuming JSON via JavaScript's Fetch API.

## Philosophy
The library should be small enough that the code and the tests can be consumed as the documentation.

## Usage

### Paths

- parameters are defined with a leading colon ":"
    - the router disallows path prefixes followed by a different parameter name. For example, registering both of these paths would lead to an error:

    "/foo/bar/:buzz"
    "/foo/bar/:bazz"

- wildcards are defined with a leading asterisk "*"

- the match precedence for a path is:
  static -> :parameter -> *wildcard

### Context Parameters

- any values read from the URL are stored in the request context
- a `map[string]string` of parameter value key-value pairs can be retrieved with `GetParams(req.Context())`
- if there are no params, expect an empty `map[string]string`
- users should check that a value exists in the map using the standard Go idiom: `val, exists := params[key]`

