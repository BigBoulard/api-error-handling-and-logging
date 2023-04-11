package repository

import (
	"github.com/BigBoulard/api-error-handling-and-logging/microservice2/src/rest_errors"
	"github.com/jackc/pgerrcode"
)

type productRepo struct{}

type ProductRepo interface {
	GetByID() (bool, rest_errors.RestErr)
}

func NewRepo() ProductRepo {
	return &productRepo{}
}

func (r *productRepo) GetByID() (bool, rest_errors.RestErr) {
	errorCode := pgerrcode.TooManyConnections
	errorMsg := "too many connections for role backend_user"
	return false, rest_errors.NewInternalServerError("productrepo.GetByID", errorCode, errorMsg)
}
