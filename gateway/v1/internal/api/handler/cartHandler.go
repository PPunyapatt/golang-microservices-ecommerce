package handler

import (
	"context"
	"gateway/v1/proto/cart"
	"time"

	"github.com/gofiber/fiber/v2"
)

func (c *ApiHandler) AddItem(ctx *fiber.Ctx) error {
	context_, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	res, err := c.CartSvc.AddItem(context_, &cart.AddItemRequest{})
	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{
			"error": "Failed to add item to cart",
		})
	}

	return ctx.Status(200).JSON(fiber.Map{
		"message": res.Status,
	})
}

func (c *ApiHandler) RemoveItem(ctx *fiber.Ctx) error {

	return ctx.Status(200).JSON(fiber.Map{
		"message": "Add item to cart",
	})
}
