package service

import (
	"github.com/BigBoulard/api-error-handling-and-logging/src/httpclient"
	"github.com/BigBoulard/api-error-handling-and-logging/src/rest_errors"
)

type Service interface {
	DoSomething() (bool, rest_errors.RestErr)
}

type service struct {
	httpClient httpclient.HttpClient
}

func NewService(httpClient httpclient.HttpClient) Service {
	return &service{
		httpClient: httpClient,
	}
}

func (s *service) DoSomething() (bool, rest_errors.RestErr) {
	boolRes, restErr := s.httpClient.DoSomething()
	if restErr != nil {
		return false, restErr
	}
	return boolRes, nil
}
