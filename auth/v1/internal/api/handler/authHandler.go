package handler

// import (
// 	"auth-service/v1/internal/constant"
// 	"auth-service/v1/internal/helper"
// 	"auth-service/v1/proto/auth"
// 	"net/http"
// 	"time"

// 	"github.com/gofiber/fiber/v2"
// 	"github.com/google/uuid"
// 	"golang.org/x/crypto/bcrypt"
// )

// type AuthHandler struct {
// 	// Svc service.AuthService
// 	Svc auth.AuthServiceServer
// }

// // func NewAuthHandler(svc service.AuthService) *AuthHandler {
// // 	return &AuthHandler{svc}
// // }

// func (api *AuthHandler) Login(ctx *fiber.Ctx) error {
// 	// Implement the login logic here

// 	// Set token in cookie
// 	req, err := helper.ParseAndValidateRequest(ctx, &constant.User{})
// 	if err != nil {
// 		return helper.ResponseHttpError(ctx, err)
// 	}

// 	token, err := api.Svc.Login(req)
// 	if err != nil {
// 		return helper.ResponseHttpError(ctx, helper.NewHttpErrorWithDetail(http.StatusInternalServerError, err))
// 	}

// 	ctx.Cookie(&fiber.Cookie{
// 		Name:     "access-token",
// 		Value:    *token,
// 		Expires:  time.Now().Add(24 * time.Hour),
// 		HTTPOnly: true,
// 		Secure:   false, // set to true in production with HTTPS
// 		Path:     "/",
// 	})

// 	return ctx.Status(http.StatusOK).JSON(constant.DataResponse{
// 		StatusCode: http.StatusOK,
// 		Data:       token,
// 		Message:    "Login successful",
// 	})
// }

// func (api *AuthHandler) Register(ctx *fiber.Ctx) error {
// 	// Implement the register logic here
// 	req, err := helper.ParseAndValidateRequest(ctx, &constant.User{})
// 	if err != nil {
// 		return helper.ResponseHttpError(ctx, err)
// 	}

// 	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
// 	if err != nil {
// 		return helper.ResponseHttpError(ctx, helper.NewHttpErrorWithDetail(http.StatusInternalServerError, err))
// 	}

// 	req.ID = uuid.NewString()
// 	req.CreatedAt = time.Now()
// 	req.Password = string(hashedPassword)

// 	err = api.Svc.Register(req)
// 	if err != nil {
// 		return helper.ResponseHttpError(ctx, helper.NewHttpErrorWithDetail(http.StatusInternalServerError, err))
// 	}

// 	return ctx.Status(http.StatusOK).JSON(constant.StatusResponse{
// 		StatusCode: http.StatusOK,
// 		Message:    "Register successful",
// 	})
// }
