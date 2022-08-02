package domain

import (
	"capi/errs"
	"capi/logger"
	"database/sql"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type CustomerRepositoryDB struct {
	db *sqlx.DB
}

func NewCustomerRepositoryDB(client *sqlx.DB) CustomerRepositoryDB {

	return CustomerRepositoryDB{client}
}

func (s CustomerRepositoryDB) FindAll(status string) ([]Customer, *errs.AppErr) {

	var query string
	var err error

	var customers []Customer

	if status == "" {
		query = "select * from customers"
		err = s.db.Select(&customers, query)
	} else {
		query = "select * from customers where status = $1"
		err = s.db.Select(&customers, query, status)
	}

	if err != nil {
		logger.Error("error fetch data to customer table " + err.Error())
		return nil, errs.NewUnexpectedError("unexpected database error")
	}

	return customers, nil

}

func (s CustomerRepositoryDB) FindByID(id string) (*Customer, *errs.AppErr) {

	query := "select * from customers where customer_id = $1"

	var c Customer

	err := s.db.Get(&c, query, id)

	if err != nil {
		if err == sql.ErrNoRows {
			logger.Error(err.Error())
			return nil, errs.NewNotFoundError("customer not found")
		} else {
			log.Println("error scanning customer data ", err.Error())
			return nil, errs.NewUnexpectedError("unexpected database error")
		}

	}

	return &c, nil

}
