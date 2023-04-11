package service

import (
	"github.com/BigBoulard/api-error-handling-and-logging/microservice2/src/domain"
	"github.com/BigBoulard/api-error-handling-and-logging/microservice2/src/httpclient"
	"github.com/BigBoulard/api-error-handling-and-logging/microservice2/src/rest_errors"
)

type Service interface {
	GetTodos() ([]domain.Todo, rest_errors.RestErr)
}

type service struct {
	todosHttpClient httpclient.TodosHttpClient
}

func NewService(todosHttpClient httpclient.TodosHttpClient) Service {
	return &service{
		todosHttpClient: todosHttpClient,
	}
}

func (s *service) GetTodos() ([]domain.Todo, rest_errors.RestErr) {
	todos, restErr := s.todosHttpClient.GetTodos()
	if restErr != nil {
		restErr.WrapPath("todosservice.GetTodos/")
		return nil, restErr
	}
	return todos, nil
}
