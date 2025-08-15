package handler

import (
	"context"
	"errors"
	"gateway/v1/internal/constant"
	"gateway/v1/internal/helper"
	"gateway/v1/proto/order"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/goforj/godump"
)

func (c *ApiHandler) PlaceOrder(ctx *fiber.Ctx) error {
	request, err := helper.ParseAndValidateRequest(ctx, &constant.PlaceOrderRequest{})
	if err != nil {
		return helper.RespondHttpError(ctx, err)
	}

	context_, cancel := context.WithTimeout(ctx.UserContext(), time.Second*10)
	defer cancel()

	userID, ok := ctx.Locals("UserID").(string)
	if !ok {
		return helper.RespondHttpError(ctx, errors.New("user ID not found in context"))
	}

	godump.Dump(request)

	orderItems := []*order.OrderItems{}

	for _, orderStore := range request.OrderItems {
		orderStores := &order.OrderItems{
			StoreId: orderStore.StoreID,
		}
		items := []*order.Item{}
		for _, item := range orderStore.Items {
			item := &order.Item{
				ProductId: item.ProductID,
				Quantity:  item.Quantity,
			}

			items = append(items, item)
		}

		orderStores.Items = items
		orderItems = append(orderItems, orderStores)
	}

	godump.Dump(orderItems)

	_, err = c.OrderSvc.PlaceOrder(context_, &order.PlaceOrderRequest{
		UserId:     userID,
		ShippingId: request.ShippingID,
		OrderItems: orderItems,
	})
	if err != nil {
		return helper.RespondHttpError(ctx, err)
	}

	return ctx.Status(200).JSON(fiber.Map{
		"status": "Success",
	})
}
