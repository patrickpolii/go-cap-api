package service

import (
	"capi/domain"
	"capi/dto"
	"capi/errs"
	"time"
)

const dbTSLayout = "2006-01-02 15:04:05"

type AccountService interface {
	NewAccount(dto.NewAccountRequest) (*dto.NewAccountResponse, *errs.AppErr)
	MakeTransaction(dto.TransactionRequest) (*dto.TransactionResponse, *errs.AppErr)
}

type DefaultAccountService struct {
	repository domain.AccountRepository
}

func NewAccountService(repo domain.AccountRepository) DefaultAccountService {
	return DefaultAccountService{repo}
}

func (s DefaultAccountService) NewAccount(req dto.NewAccountRequest) (*dto.NewAccountResponse, *errs.AppErr) {

	if err := req.Validate(); err != nil {
		return nil, err
	}

	account := domain.NewAccount(req.CustomerID, req.AccountType, req.Amount)

	newAcc, err := s.repository.Save(account)
	if err != nil {
		return nil, err
	}

	response := newAcc.ToNewAccountResponseDTO()

	return &response, nil
}

func (s DefaultAccountService) MakeTransaction(req dto.TransactionRequest) (*dto.TransactionResponse, *errs.AppErr) {
	// incoming request validation
	err := req.Validate()
	if err != nil {
		return nil, err
	}
	// server side validation for checking the available balance in the account
	if req.IsTransactionTypeWithdrawal() {
		account, err := s.repository.FindBy(req.AccountID)
		if err != nil {
			return nil, err
		}
		if account.CanWithdraw(req.Amount) {
			return nil, errs.NewValidationError("Insufficient balance in the account")
		}
	}
	// if all is well, build the domain object & save the transaction
	t := domain.Transaction{
		AccountId:       req.AccountID,
		Amount:          req.Amount,
		TransactionType: req.TransactionType,
		TransactionDate: time.Now().Format(dbTSLayout),
	}
	transaction, appError := s.repository.SaveTransaction(t)
	if appError != nil {
		return nil, appError
	}
	response := transaction.ToDto()
	return &response, nil
}
