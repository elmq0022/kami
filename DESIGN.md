# Requirements

1. No dependencies; standard library only.


## Methods
Support all standard HTTP methods, GET, POST, PUT, DELETE

## URL matching


## Error Handling
Handle panics gracefully.
A panic should not kill the http server.

An error should return a status code in the 4xx - 5xx range.
Errors should inclued a JSON payload consisting of:
- error
- code
- message
- details

## Content Type
Always Content-Type: application/json

## Ergonomics
- shouldn't have to marshall/unmarshall all the time
- there's a good error struct(s) provided
