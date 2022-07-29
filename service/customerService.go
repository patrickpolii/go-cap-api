package service

import (
	"capi/domain"
	"capi/dto"
	"capi/errs"
)

type CustomerService interface {
	GetAllCustomer(string) ([]dto.CustomerResponse, *errs.AppErr)
	GetCustomerByID(string) (*dto.CustomerResponse, *errs.AppErr)
}

type DefaultCustomerService struct {
	repository domain.CustomerRepository
}

func (s DefaultCustomerService) GetAllCustomer(customerStatus string) ([]dto.CustomerResponse, *errs.AppErr) {
	// * add process here
	customers, err := s.repository.FindAll(customerStatus)
	if err != nil {
		return nil, errs.NewUnexpectedError("unexpected db error")
	}
	var dtoCustomers []dto.CustomerResponse
	for _, customer := range customers {
		dtoCustomers = append(dtoCustomers, customer.ToDTO())
	}
	return dtoCustomers, nil

	//return s.repository.FindAll()
}

func (s DefaultCustomerService) GetCustomerByID(customerID string) (*dto.CustomerResponse, *errs.AppErr) {
	cust, err := s.repository.FindByID(customerID)
	if err != nil {
		return nil, err
	}

	response := cust.ToDTO()

	return &response, nil
}

func NewCustomerService(repository domain.CustomerRepository) DefaultCustomerService {
	return DefaultCustomerService{repository}
}
