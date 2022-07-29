package domain

import (
	"capi/dto"
	"capi/errs"
)

type Customer struct {
	ID          string `json:"id" xml:"id" db:"customer_id"`
	Name        string `json:"name" xml:"name"`
	City        string `json:"city" xml:"city"`
	ZipCode     string `json:"zip_code" xml:"zipcode"`
	DateOfBirth string `json:"date_of_birth" xml:"dateofbirth" db:"date_of_birth"`
	Status      string `json:"status" xml:"status"`
}

type CustomerRepository interface {
	FindAll(string) ([]Customer, *errs.AppErr)
	FindByID(string) (*Customer, *errs.AppErr)
}

func (c Customer) convertStatusName() string {
	statusName := "active"
	if c.Status == "0" {
		statusName = "inactive"
	}
	return statusName
}

func (c Customer) ToDTO() dto.CustomerResponse {
	return dto.CustomerResponse{
		ID:          c.ID,
		Name:        c.Name,
		DateOfBirth: c.DateOfBirth,
		City:        c.City,
		ZipCode:     c.ZipCode,
		Status:      c.convertStatusName(),
	}
}
