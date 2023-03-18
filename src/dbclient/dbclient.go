package dbclient

import (
	"errors"

	"github.com/BigBoulard/api-error-handling-and-logging/src/rest_errors"
)

type client struct{}

type DBClient interface {
	DoSomething() (bool, rest_errors.RestErr)
}

func NewClient() DBClient {
	return &client{}
}

func (c *client) DoSomething() (bool, rest_errors.RestErr) {
	noRecordError := errors.New("the resource doesn't exist")
	return false, rest_errors.NewNotFoundError(noRecordError.Error())
}
