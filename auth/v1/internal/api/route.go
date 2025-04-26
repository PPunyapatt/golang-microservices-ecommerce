package api

import (
	"auth-service/v1/internal/api/handler"
	"auth-service/v1/internal/service"

	"github.com/gofiber/fiber/v2"
)

// func Route(app *fiber.App, api handler.AuthHandler) {
func Route(app *fiber.App, svc service.AuthService) {
	api := &handler.AuthHandler{Svc: svc}

	auth := app.Group("/api/v1/auth")
	auth.Post("/login", api.Login)
	auth.Post("/register", api.Register)
}
