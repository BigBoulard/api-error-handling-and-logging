package productsrepo

import (
	"github.com/BigBoulard/api-error-handling-and-logging/microservice1/src/rest_errors"
	"github.com/jackc/pgerrcode"
)

type repo struct{}

type Repo interface {
	GetByID() (bool, rest_errors.RestErr)
}

func NewRepo() Repo {
	return &repo{}
}

func (r *repo) GetByID() (bool, rest_errors.RestErr) {
	return false, rest_errors.NewInternalServerError(
		"productrepo.GetByID",
		pgerrcode.TooManyConnections,
		"too many connections for role backend_user",
	)
}
