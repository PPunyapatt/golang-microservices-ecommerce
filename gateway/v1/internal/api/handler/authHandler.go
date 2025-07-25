package handler

import (
	"context"
	"gateway/v1/internal/constant"
	"gateway/v1/proto/auth"
	"time"

	"github.com/gofiber/fiber/v2"
)

func (c *ApiHandler) Login(ctx *fiber.Ctx) error {
	// Implement login logic here
	user := &constant.User{}
	if err := ctx.BodyParser(user); err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	context_, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	res, err := c.AuthSvc.Login(context_, &auth.LoginRequest{
		Email:    user.Email,
		Password: user.Pasword,
	})
	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{
			"error": "Failed to login: " + err.Error(),
		})
	}
	return ctx.Status(200).JSON(fiber.Map{
		"message": res.Token,
	})
}

func (c *ApiHandler) Register(ctx *fiber.Ctx) error {
	user := &constant.UserRegister{}
	if err := ctx.BodyParser(user); err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	context_, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	res, err := c.AuthSvc.Register(context_, &auth.RegisterRequest{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Password:  user.Password,
		Email:     user.Email,
		Phone:     user.Phone,
	})

	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{
			"error": "Failed to login" + err.Error(),
		})
	}

	return ctx.Status(200).JSON(fiber.Map{
		"message": res.Status,
	})
}

func (c *ApiHandler) CreateSrore(ctx *fiber.Ctx) error {
	store := &constant.Store{}
	if err := ctx.BodyParser(store); err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	context_, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	userID, ok := ctx.Locals("UserID").(string)
	if !ok {
		return ctx.Status(500).JSON(fiber.Map{
			"error": "Don't have UserID",
		})
	}

	res, err := c.AuthSvc.CreateStore(context_, &auth.CreateStoreRequest{
		Name:   store.Name,
		UserID: userID,
	})

	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{
			"error": "Create failed: " + err.Error(),
		})
	}
	return ctx.Status(200).JSON(fiber.Map{
		"message": res.Status,
	})
}
