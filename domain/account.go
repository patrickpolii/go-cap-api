package domain

import (
	"capi/dto"
	"capi/errs"
	"time"
)

const dbTSLayout = "2006-01-02 15:04:05"

type Account struct {
	AccountID   string  `db:"account_id"`
	CustomerID  string  `db:"customer_id"`
	OpeningDate string  `db:"opening_date"`
	AccountType string  `db:"account_type"`
	Amount      float64 `db:"amount"`
	Status      string  `db:"status"`
}

func (a Account) ToNewAccountResponseDTO() dto.NewAccountResponse {
	return dto.NewAccountResponse{AccountID: a.AccountID}
}

type AccountRepository interface {
	Save(Account) (*Account, *errs.AppErr)
	SaveTransaction(transaction Transaction) (*Transaction, *errs.AppErr)
	FindBy(accountId string) (*Account, *errs.AppErr)
}

func (a Account) CanWithdraw(amount float64) bool {
	return a.Amount < amount
}

func NewAccount(customerId, accountType string, amount float64) Account {
	return Account{
		CustomerID:  customerId,
		OpeningDate: time.Now().Format(dbTSLayout),
		AccountType: accountType,
		Amount:      amount,
		Status:      "1",
	}
}
