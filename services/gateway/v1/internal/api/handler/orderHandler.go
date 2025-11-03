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
	c.prometheusMetrics.Http.OrderPlaceRequests.Inc()
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

	orderStores := []*order.PlaceOrderStores{}

	for _, orderStore := range request.OrderItems {
		orderStore_ := &order.PlaceOrderStores{
			StoreId: orderStore.StoreID,
		}
		items := []*order.PlaceOrderItems{}
		for _, item := range orderStore.Items {
			item := &order.PlaceOrderItems{
				ProductId: item.ProductID,
				Quantity:  item.Quantity,
			}

			items = append(items, item)
		}

		orderStore_.Items = items
		// orderStores = append(orderStores, orderStores)
		orderStores = append(orderStores, orderStore_)
	}

	godump.Dump(orderStores)

	_, err = c.OrderSvc.PlaceOrder(context_, &order.PlaceOrderRequest{
		UserId:      userID,
		ShippingId:  request.ShippingID,
		OrderItems:  orderStores,
		OrderSource: request.OrderSource,
	})
	if err != nil {
		return helper.RespondHttpError(ctx, err)
	}

	return ctx.Status(200).JSON(fiber.Map{
		"status": "Success",
	})
}

func (c *ApiHandler) ListOrder(ctx *fiber.Ctx) error {
	status := ctx.Params("status")
	context_, cancel := context.WithTimeout(ctx.UserContext(), time.Second*10)
	defer cancel()

	userID, ok := ctx.Locals("UserID").(string)
	if !ok {
		return helper.RespondHttpError(ctx, errors.New("user ID not found in context"))
	}

	_, err := c.OrderSvc.ListOrder(context_, &order.ListOrderRequest{
		UserId: userID,
		Status: &status,
	})
	if err != nil {
		return helper.RespondHttpError(ctx, err)
	}

	return ctx.Status(200).JSON(fiber.Map{
		"status": "Success",
	})
}
