package service

import (
	"auth-service/v1/internal/repository"
	"auth-service/v1/proto/auth"
	"context"
)

type authServer struct {
	authRepo repository.AuthRepository
	auth.UnimplementedAuthServiceServer
}

func NewAuthServer(authRepo repository.AuthRepository) auth.AuthServiceServer {
	return authServer{authRepo: authRepo}
}

func (s authServer) Login(ctx context.Context, in *auth.LoginRequest) (*auth.LoginResponse, error) {

	loginResponse := auth.LoginResponse{
		Token: "mocked-jwt-token",
	}

	return &loginResponse, nil
}

func (s authServer) Register(ctx context.Context, in *auth.RegisterLogin) (*auth.RegisterResponse, error) {
	registerResponse := auth.RegisterResponse{
		Status: "success",
	}

	return &registerResponse, nil
}
