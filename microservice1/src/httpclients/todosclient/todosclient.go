package todosclient

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/BigBoulard/api-error-handling-and-logging/microservice1/src/conf"
	"github.com/BigBoulard/api-error-handling-and-logging/microservice1/src/domains/todosdom"

	"github.com/BigBoulard/api-error-handling-and-logging/microservice1/src/rest_errors"
	"github.com/go-resty/resty/v2"
)

type client struct {
	resty *resty.Client
}

type Client interface {
	GetTodos() ([]todosdom.Todo, rest_errors.RestErr)
}

func NewClient() Client {
	conf.LoadEnv() // load env vars
	r := resty.New()
	if conf.Env.AppMode == "debug" { // Debug mode is taken from the env vars
		r.SetDebug(true)
	}
	r.SetBaseURL(fmt.Sprintf("http://%s:%s", conf.Env.Microservice2Host, conf.Env.Microservice2Port))

	return &client{
		resty: r,
	}
}

func (c *client) GetTodos() ([]todosdom.Todo, rest_errors.RestErr) {
	resp, err := c.resty.
		R().
		SetHeader("Accept", "application/json").
		Get("/todos")

	if err != nil {
		return nil, rest_errors.NewInternalServerError("httpclient/GetTodos", "", err.Error())
	}

	if resp.StatusCode() > 399 {
		return nil, rest_errors.NewInternalServerError("httpclient/GetTodos", strconv.Itoa(resp.StatusCode()), err.Error())
	}

	var todos []todosdom.Todo
	err = json.Unmarshal(resp.Body(), &todos)
	if err != nil {
		return nil, rest_errors.NewInternalServerError("httpclient/GetTodos", "", err.Error())
	}

	return todos, nil
}
