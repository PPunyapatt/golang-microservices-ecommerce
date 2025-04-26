package service

import (
	"auth-service/v1/internal/constant"
	"auth-service/v1/internal/helper"
	"auth-service/v1/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Login(*constant.User) (*string, error)
	Register(*constant.User) error
}

type authService struct {
	authRepo repository.AuthRepository
}

func NewAuthService(authRepo repository.AuthRepository) AuthService {
	return &authService{
		authRepo: authRepo,
	}
}

func (service *authService) Login(user *constant.User) (*string, error) {
	u, err := service.authRepo.Login(user)
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(user.Password)); err != nil {
		return nil, err //errors.New("password is invalid")
	}

	token, err := helper.GenerateToken(u.ID, u.Role)
	if err != nil {
		return nil, nil
	}

	return token, nil
}

func (service *authService) Register(user *constant.User) error {
	err := service.authRepo.Register(user)
	if err != nil {
		return err
	}
	return nil
}
