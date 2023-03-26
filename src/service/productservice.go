package service

import (
	"github.com/BigBoulard/api-error-handling-and-logging/src/repository"
	"github.com/BigBoulard/api-error-handling-and-logging/src/rest_errors"
)

type Service interface {
	GetProduct() (bool, rest_errors.RestErr)
}

type service struct {
	productRepo repository.ProductRepo
}

func NewService(productRepo repository.ProductRepo) Service {
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
