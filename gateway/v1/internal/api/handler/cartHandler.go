package handler

import (
	"context"
	"errors"
	"gateway/v1/internal/constant"
	"gateway/v1/internal/helper"
	"gateway/v1/proto/cart"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/goforj/godump"
)

func (c *ApiHandler) AddItem(ctx *fiber.Ctx) error {
	request, err := helper.ParseAndValidateRequest(ctx, &constant.Products{})
	if err != nil {
		return helper.RespondHttpError(ctx, err)
	}

	userID, ok := ctx.Locals("UserID").(string)
	if !ok {
		return helper.RespondHttpError(ctx, errors.New("user ID not found in context"))
	}

	godump.Dump(request)

	cartItems := []*cart.CartItem{}
	for _, item := range request.Products {
		cartItems = append(cartItems, &cart.CartItem{
			ProductId:   item.ProductID,
			ProductName: item.ProductName,
			Quantity:    item.Quantity,
			Price:       item.Price,
			ImageUrl:    item.ImageURL,
			StoreID:     item.StoreID,
		})
	}

	context_, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	res, err := c.CartSvc.AddItemToCart(context_, &cart.AddItemRequest{
		UserId: userID,
		Items:  cartItems,
	})
	if err != nil {
		return helper.RespondHttpError(ctx, err)
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
