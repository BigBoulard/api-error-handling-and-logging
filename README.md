# Private API error handling and logging

## Description

`microservice1` is a simplified example of a microservice that responds to a browser (client).
`microservice1` exposes 2 endpoints:

1. `/products/001` that generates an error. The goal here is double:
   1. check that the error returned to the client is REST friendly and contains enough info for the frontend to process the error.
   2. chech that the log contains enough information to be able to troubleshoot if necessary.
2. `/todos` that is calling microservice2 which will call an external dummy API and return a TODO list. This will be used to test concurrency later on.

## Test

1. Start microservices

```sh
> cd microservice1
> go run src/main.go
```

in another terminal

```sh
> cd microservice2
> go run src/main.go
```

2. test the microservice1 endpoints in `test.rest` using the vscode REST client extension <https://marketplace.visualstudio.com/items?itemName=humao.rest-client>

## Goals

If an error happens, the client should be able to provide an appropriate error message to the end user based on the payload received by this server
without leaking any sensitive data or giving any information about the backend code structure or the underlying technology stack.

Error Logging should provide enough information to understand the origin of an error.

## Assumptions (to be challenged)

1 - An error must be created upon detection, in a controller, service, repository or http client
2 - An error should be only handled once, and this should happen in the controller
3 - Logging should only happens in the controller and should provide the path of the request that generated the error.

## Proposal

Any error will be encapsulated into a `restErr` struct.
The `restErr` has 2 purposes:

1. Provide fields to create an standardized HTTP Error Response:

- the **HTTTP Status**,
- the **title** of the error
- a potential **cause** (like the name of a missing/invalid field).
  These fields are marshalled and presented to the client.

2. Provide Logging information to a logger: these data may contain information sensitive information (db field names, db sytem name etc.) and therefore are not marshalled.

- `ErrPath`: contains something like `productcontroller.GetProduct/productservice.GetProduct/productrepo.GetByID` to inform that an error has been triggered by the `GetByID` function from the **product repository**, called by the `GetProduct` function from the **product service** etc.
- `ErrCode` and `ErrMessage`: the raw error code and message retrieved from MySQL, PostGres, an external API or another microservice.

```go
type restErr struct {
	// these fields are returned to the browser
	ErrStatus  int    `json:"status"`          // HTTP Status Code
	ErrTitle   string `json:"title"`           // A string representation of the Status Code
	ErrCause   string `json:"cause,omitempty"` // The cause of the error, can be empty

 // these fiedls are NOT returned to the browser and are only used for logging
	ErrPath    string `json:"-"`               // The path of the error. Ex: "controller/controllerfunc/service/servicefunc/dbclient/dblientfunc"
	ErrMessage string `json:"-"`               // Raw error message returned by a DB, another Servive or whatever
	ErrCode    string `json:"-"`               // Raw error code from the DB or another service
}

type RestErr interface {
	Status() int     // HTTP status code
	Title() string   // A string representation of the Status Code

	Path() string    // The path of the error. Ex: "controller/controllerfunc/service/servicefunc/dbclient/dblientfunc"
	WrapPath(string) // Wrapper func to keep track of the path of the error
	Code() string    // Raw error code
	Message() string // Raw error message not returned to the client

	Error() string   // string representation of a restErr
}
```
