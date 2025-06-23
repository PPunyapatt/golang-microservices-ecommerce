package handler

import (
	"context"
	"gateway/v1/proto/auth"
	"time"

	"github.com/gofiber/fiber/v2"
)

type User struct {
	Email   string `json:"email"`
	Pasword string `json:"password"`
}

func (c *ApiHandler) Login(ctx *fiber.Ctx) error {
	// Implement login logic here
	p := new(User)
	if err := ctx.BodyParser(p); err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	context_, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	res, err := c.AuthSvc.Login(context_, &auth.LoginRequest{
		Email:    p.Email,
		Password: p.Pasword,
	})
	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{
			"error": "Failed to login",
		})
	}
	return ctx.Status(200).JSON(fiber.Map{
		"message": res.Token,
	})
}
