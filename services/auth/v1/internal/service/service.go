package service

import (
	"auth-service/v1/internal/constant"
	"auth-service/v1/internal/helper"
	"auth-service/v1/internal/repository"
	"auth-service/v1/proto/auth"
	"context"
	"log"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/crypto/bcrypt"
)

type authServer struct {
	tracer   trace.Tracer
	authRepo repository.AuthRepository
	auth.UnimplementedAuthServiceServer
}

func NewAuthServer(authRepo repository.AuthRepository, tracer trace.Tracer) auth.AuthServiceServer {
	return &authServer{
		authRepo: authRepo,
		tracer:   tracer,
	}
}

func (s *authServer) Login(ctx context.Context, in *auth.LoginRequest) (*auth.LoginResponse, error) {
	loginCtx, loginSpan := s.tracer.Start(ctx, "Login")
	user := &constant.User{
		Email:    in.Email,
		Password: in.Password,
	}

	u, err := s.authRepo.Login(loginCtx, user)
	if err != nil {
		return nil, err
	}

	_, compareSpan := s.tracer.Start(ctx, "Compare Password")
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(user.Password)); err != nil {
		log.Println("compare: ", err.Error())
		return nil, err //errors.New("password is invalid")
	}
	compareSpan.End()

	_, tokenSpan := s.tracer.Start(ctx, "GenerateToken")
	token, err := helper.GenerateToken(u.ID, u.Roles)
	if err != nil {
		return nil, nil
	}
	tokenSpan.End()

	loginResponse := auth.LoginResponse{
		Token: *token,
	}
	loginSpan.End()

	return &loginResponse, nil
}

func (s *authServer) Register(ctx context.Context, in *auth.RegisterRequest) (*auth.RegisterResponse, error) {
	regisCtx, regisSpan := s.tracer.Start(ctx, "Register")

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
	err = s.authRepo.Register(regisCtx, user)
	if err != nil {
		return nil, err
	}
	registerResponse := auth.RegisterResponse{
		Status: "success",
	}

	regisSpan.End()

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
