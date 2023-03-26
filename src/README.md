# Private API error handling and logging

## Test

```go
go run src/main.go
```

## Context

This code is a simplified example of a private API that responds to a browser (client).

## Goals

If an error happens, the client should be able to provide an appropriate error message to the end user based on the payload received by this server
without leaking any sensitive data or giving any information about the backend code structure or the underlying technology stack.

Error Logging should provide enough information to understand the origin of an error:

## Assumptions (to be challenged)

1 - An error must be created upon detection, in a controller, service, repository or http client
2 - On the other hand, logging should only happens in the controller and should provide enough information to understand the path of the request that generated the error.

## Proposal

A `restErr` struct is created and contains:

1. Necessary fields to give a correct HTTP Error Response:

- the **HTTTP Status**,
- the **title** of the error
- a potential **cause** (like the name of a missing/invalid field).
  These fields are marshalled and presented to the client.

2. Logging information: these data may contain information sensitive information (db field names, db sytem name etc.) and therefore are not marshalled.

- `ErrPath`: contains something like `productcontroller.GetProduct/productservice.GetProduct/productrepo.GetByID` to inform that an error has been triggered by the `GetByID` function from the **product repository**, called by the `GetProduct` function from the **product service** etc.
- `ErrCode` and `ErrMessage`: the raw error code and message retrieved from MySQL, PostGres, an external API or another microservice.

```go
type restErr struct {
	ErrStatus  int    `json:"status"`          // HTTP Status Code
	ErrTitle   string `json:"title"`           // A string representation of the Status Code
	ErrCause   string `json:"cause,omitempty"` // The cause of the error, can be empty

	ErrPath    string `json:"-"`               // Only used for Logging: The path of the error. Ex: "controller/controllerfunc/service/servicefunc/dbclient/dblientfunc"
	ErrMessage string `json:"-"`               // Only used for Logging: Raw error message returned by a DB, another Servive or whatever
	ErrCode    string `json:"-"`               // Only used for Logging: Raw error code from the DB or another service
}

type RestErr interface {
	Status() int     // HTTP status code
	Title() string   // A string representation of the Status Code

	Path() string    // Only used for Logging: The path of the error. Ex: "controller/controllerfunc/service/servicefunc/dbclient/dblientfunc"
	WrapPath(string) // Only used for Logging: Wrapper func to keep track of the path of the error
	Code() string    // Only used for Logging: Raw error code
	Message() string // Only used for Logging: Raw error message not returned to the client

	Error() string   // string representation of a restErr
}
```
