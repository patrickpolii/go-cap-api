package domain

import (
	"capi/errs"
	"capi/logger"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type AccountRepositoryDB struct {
	db *sqlx.DB
}

func NewAccountRepositoryDB(client *sqlx.DB) AccountRepositoryDB {
	return AccountRepositoryDB{client}
}

func (d AccountRepositoryDB) Save(a Account) (*Account, *errs.AppErr) {
	query := "insert into accounts (customer_id, opening_date, account_type, amount, status) values ($1, $2, $3, $4, $5) RETURNING account_id"

	var id string
	err := d.db.QueryRow(query, a.CustomerID, a.OpeningDate, a.AccountType, a.Amount, a.Status).Scan(&id)
	if err != nil {
		logger.Error(fmt.Sprintf("error while creating new account: %v "+err.Error(), a))
		return nil, errs.NewUnexpectedError("unexpected error from database")
	}

	a.AccountID = id
	return &a, nil
}

/**
 * transaction = make an entry in the transaction table + update the balance in the accounts table
 */
func (d AccountRepositoryDB) SaveTransaction(t Transaction) (*Transaction, *errs.AppErr) {
	// starting the database transaction block
	tx, err := d.db.Begin()
	if err != nil {
		logger.Error("Error while starting a new transaction for bank account transaction: " + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected database error")
	}

	// inserting bank account transaction
	var id string
	err = tx.QueryRow(`INSERT INTO transactions (account_id, amount, transaction_type, transaction_date) 
						values ($1, $2, $3, $4) RETURNING transaction_id`, t.AccountId, t.Amount, t.TransactionType, t.TransactionDate).Scan(&id)
	if err != nil {
		logger.Error("Error while starting a new transaction for bank account transaction: " + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected database error")
	}

	// updating account balance
	if t.IsWithdrawal() {
		_, err = tx.Exec(`UPDATE accounts SET amount = amount - $1 where account_id = $2`, t.Amount, t.AccountId)
	} else {
		_, err = tx.Exec(`UPDATE accounts SET amount = amount + $1 where account_id = $2`, t.Amount, t.AccountId)
	}

	// in case of error Rollback, and changes from both the tables will be reverted
	if err != nil {
		tx.Rollback()
		logger.Error("Error while saving transaction: " + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected database error")
	}
	// commit the transaction when all is good
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		logger.Error("Error while commiting transaction for bank account: " + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected database error")
	}

	// Getting the latest account information from the accounts table
	account, appErr := d.FindBy(t.AccountId)
	if appErr != nil {
		return nil, appErr
	}
	t.TransactionId = id

	// updating the transaction struct with the latest balance
	t.Amount = account.Amount
	return &t, nil
}

func (d AccountRepositoryDB) FindBy(accountId string) (*Account, *errs.AppErr) {
	sqlGetAccount := "SELECT account_id, customer_id, opening_date, account_type, amount from accounts where account_id = $1"
	var account Account
	err := d.db.Get(&account, sqlGetAccount, accountId)
	if err != nil {
		logger.Error("Error while fetching account information: " + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected database error")
	}
	return &account, nil
}
