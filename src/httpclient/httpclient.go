package httpclient

import (
	"github.com/BigBoulard/api-error-handling-and-logging/src/rest_errors"
)

type client struct{}

type HttpClient interface {
	DoSomething() (bool, rest_errors.RestErr)
}

func NewClient() HttpClient {
	return &client{}
}

func (c *client) DoSomething() (bool, rest_errors.RestErr) {

	return true, nil
}
