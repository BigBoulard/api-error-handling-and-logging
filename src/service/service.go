package service

import (
	"github.com/BigBoulard/api-error-handling-and-logging/src/dbclient"
	"github.com/BigBoulard/api-error-handling-and-logging/src/rest_errors"
)

type Service interface {
	DoSomething() (bool, rest_errors.RestErr)
}

type service struct {
	dbClient dbclient.DBClient
}

func NewService(dbClient dbclient.DBClient) Service {
	return &service{
		dbClient: dbClient,
	}
}

func (s *service) DoSomething() (bool, rest_errors.RestErr) {
	boolRes, restErr := s.dbClient.DoSomething()
	if restErr != nil {
		return false, restErr
	}
	return boolRes, nil
}
