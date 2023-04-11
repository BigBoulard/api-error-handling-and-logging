package productsservice

import (
	"github.com/BigBoulard/api-error-handling-and-logging/microservice1/src/repositories/productsrepo"
	"github.com/BigBoulard/api-error-handling-and-logging/microservice1/src/rest_errors"
)

type Service interface {
	GetProduct() (bool, rest_errors.RestErr)
}

type service struct {
	productRepo productsrepo.Repo
}

func NewService(productRepo productsrepo.Repo) Service {
	return &service{
		productRepo: productRepo,
	}
}

func (s *service) GetProduct() (bool, rest_errors.RestErr) {
	boolRes, restErr := s.productRepo.GetByID()
	if restErr != nil {
		restErr.WrapPath("productservice.GetProduct/")
		return false, restErr
	}
	return boolRes, nil
}
