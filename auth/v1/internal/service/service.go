package service

import (
	"auth-service/v1/internal/constant"
	"auth-service/v1/internal/helper"
	"auth-service/v1/internal/repository"
	"auth-service/v1/proto/auth"
	"context"
	"log"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type authServer struct {
	authRepo repository.AuthRepository
	auth.UnimplementedAuthServiceServer
}

func NewAuthServer(authRepo repository.AuthRepository) auth.AuthServiceServer {
	return &authServer{authRepo: authRepo}
}

func (s *authServer) Login(ctx context.Context, in *auth.LoginRequest) (*auth.LoginResponse, error) {
	user := &constant.User{
		Email:    in.Email,
		Password: in.Password,
	}

	u, err := s.authRepo.Login(user)
	if err != nil {
		return nil, err
	}

	log.Println(u.Password, user.Password)

	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(user.Password)); err != nil {
		return nil, err //errors.New("password is invalid")
	}

	token, err := helper.GenerateToken(u.ID, u.Roles)
	if err != nil {
		return nil, nil
	}

	loginResponse := auth.LoginResponse{
		Token: *token,
	}

	return &loginResponse, nil
}

func (s *authServer) Register(ctx context.Context, in *auth.RegisterRequest) (*auth.RegisterResponse, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(in.Password), 14)
	if err != nil {
		// handle error
		log.Println("Error hashing password:", err)
		return nil, err
	}
	user := &constant.User{
		ID:          uuid.NewString(),
		FirstName:   in.FirstName,
		LastName:    in.LastName,
		Email:       in.Email,
		Password:    string(passwordHash),
		PhoneNumber: in.Phone,
	}
	err = s.authRepo.Register(user)
	if err != nil {
		return nil, err
	}
	registerResponse := auth.RegisterResponse{
		Status: "success",
	}

	return &registerResponse, nil
}

func (s *authServer) CreateStore(ctx context.Context, in *auth.CreateStoreRequest) (*auth.CreateStoreResponse, error) {
	store := &constant.Store{
		Name:   in.Name,
		UserID: in.UserID,
	}
	err := s.authRepo.CreateStore(store)
	if err != nil {
		return nil, err
	}

	response := auth.CreateStoreResponse{
		Status: "Create store success",
	}

	return &response, nil
}
