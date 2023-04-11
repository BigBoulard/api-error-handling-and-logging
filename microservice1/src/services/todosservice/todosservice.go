package todosservice

import (
	"github.com/BigBoulard/api-error-handling-and-logging/microservice1/src/domains/todosdom"
	"github.com/BigBoulard/api-error-handling-and-logging/microservice1/src/httpclients/todosclient"
	"github.com/BigBoulard/api-error-handling-and-logging/microservice1/src/rest_errors"
)

type Service interface {
	GetTodos() ([]todosdom.Todo, rest_errors.RestErr)
}

type service struct {
	todosClient todosclient.Client
}

func NewService(todosClient todosclient.Client) Service {
	return &service{
		todosClient: todosClient,
	}
}

func (s *service) GetTodos() ([]todosdom.Todo, rest_errors.RestErr) {
	todos, restErr := s.todosClient.GetTodos()
	if restErr != nil {
		restErr.WrapPath("productservice.GetTodos/")
		return nil, restErr
	}
	return todos, nil
}
