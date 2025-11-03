package handler

import (
	"context"
	"fmt"
	"gateway/v1/internal/constant"
	"gateway/v1/internal/helper"
	"gateway/v1/proto/auth"
	"net/http"
	"time"

	"github.com/pkg/errors"

	"github.com/gofiber/fiber/v2"
)

func (c *ApiHandler) Login(ctx *fiber.Ctx) error {
	c.prometheusMetrics.Http.AuthLoginRequests.Inc()
	// Implement login logic here
	user, err := helper.ParseAndValidateRequest(ctx, &constant.User{})
	if err != nil {
		return helper.RespondHttpError(ctx, helper.NewHttpError(http.StatusBadRequest, err))
	}

	context_, cancel := context.WithTimeout(ctx.UserContext(), time.Second*10)
	defer cancel()

	res, err := c.AuthSvc.Login(context_, &auth.LoginRequest{
		Email:    user.Email,
		Password: user.Pasword,
	})
	if err != nil {
		return helper.RespondHttpError(ctx, helper.NewHttpError(http.StatusUnauthorized, err))
	}

	return ctx.Status(200).JSON(fiber.Map{
		"message": res.Token,
	})
}

func (c *ApiHandler) Register(ctx *fiber.Ctx) error {
	c.prometheusMetrics.Http.AuthRegisterRequests.Inc()
	user := &constant.UserRegister{}
	if err := ctx.BodyParser(user); err != nil {
		return helper.RespondHttpError(ctx, helper.NewHttpError(http.StatusInternalServerError, err))
	}

	context_, cancel := context.WithTimeout(ctx.Context(), time.Second*10)
	defer cancel()

	res, err := c.AuthSvc.Register(context_, &auth.RegisterRequest{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Password:  user.Password,
		Email:     user.Email,
		Phone:     user.Phone,
	})

	if err != nil {
		return helper.RespondHttpError(ctx, err)
	}

	return ctx.Status(200).JSON(fiber.Map{
		"message": fmt.Sprintf("Register %s", res.Status),
	})
}

func (c *ApiHandler) CreateSrore(ctx *fiber.Ctx) error {
	store := &constant.Store{}
	if err := ctx.BodyParser(store); err != nil {
		return helper.RespondHttpError(ctx, helper.NewHttpError(http.StatusBadRequest, errors.Wrap(err, "Invalid request body")))
	}

	context_, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	userID, ok := ctx.Locals("UserID").(string)
	if !ok {
		return helper.RespondHttpError(ctx, errors.New("Don't have UserID"))
	}

	res, err := c.AuthSvc.CreateStore(context_, &auth.CreateStoreRequest{
		Name:   store.Name,
		UserID: userID,
	})

	if err != nil {
		return helper.RespondHttpError(ctx, errors.Wrap(err, "Cteate store failed"))
	}
	return ctx.Status(200).JSON(fiber.Map{
		"message": res.Status,
	})
}
