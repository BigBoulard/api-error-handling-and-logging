package httpclient

import (
	"encoding/json"
	"strconv"

	"github.com/BigBoulard/api-error-handling-and-logging/microservice2/src/conf"
	"github.com/BigBoulard/api-error-handling-and-logging/microservice2/src/domain"

	"github.com/BigBoulard/api-error-handling-and-logging/microservice2/src/rest_errors"
	"github.com/go-resty/resty/v2"
)

type client struct {
	Resty *resty.Client
}

type TodosHttpClient interface {
	GetTodos() ([]domain.Todo, rest_errors.RestErr)
}

func NewTodosClient() TodosHttpClient {
	conf.LoadEnv() // load env vars
	r := resty.New()
	if conf.Env.AppMode == "debug" { // Debug mode is taken from the env vars
		r.SetDebug(true)
	}
	r.SetBaseURL("http://jsonplaceholder.typicode.com")

	return &client{
		Resty: r,
	}
}

func (c *client) GetTodos() ([]domain.Todo, rest_errors.RestErr) {
	resp, err := c.Resty.
		R().
		SetHeader("Accept", "application/json").
		Get("/todos")

	if err != nil {
		return nil, rest_errors.NewInternalServerError("httpclient/GetTodos", "", err.Error())
	}

	if resp.StatusCode() > 399 {
		return nil, rest_errors.NewInternalServerError("httpclient/GetTodos", strconv.Itoa(resp.StatusCode()), err.Error())
	}

	var todos []domain.Todo
	err = json.Unmarshal(resp.Body(), &todos)
	if err != nil {
		return nil, rest_errors.NewInternalServerError("httpclient/GetTodos", "", err.Error())
	}

	return todos, nil
}
