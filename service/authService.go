package service

import (
	"capi/domain"
	"capi/dto"
	"capi/errs"
	"capi/logger"
	"strings"

	"github.com/golang-jwt/jwt/v4"
)

type AuthService interface {
	Login(dto.LoginRequest) (*dto.LoginResponse, *errs.AppErr)
}

type DefaultAuthService struct {
	repository domain.AuthRepository
}

func NewAuthService(repository domain.AuthRepository) DefaultAuthService {
	return DefaultAuthService{repository}
}

func (s DefaultAuthService) Login(req dto.LoginRequest) (*dto.LoginResponse, *errs.AppErr) {

	login, errApp := s.repository.FindBy(req.Username, req.Password)
	if errApp != nil {
		return nil, errApp
	}

	accounts := strings.Split(login.Accounts.String, ",")
	claims := domain.AccessTokenClaims{
		CustomerID: login.CustomerID.String,
		Username:   login.Username,
		Role:       login.Role,
		Accounts:   accounts,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte("rahasia"))
	if err != nil {
		logger.Error("failed while signing access token: " + err.Error())
		return nil, errs.NewUnexpectedError("cannot generate access token")
	}

	return &dto.LoginResponse{AccessToken: signedToken}, nil

}
