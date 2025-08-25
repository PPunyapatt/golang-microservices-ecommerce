package handler

import (
	"context"
	"errors"
	"gateway/v1/internal/constant"
	"gateway/v1/internal/helper"
	"gateway/v1/proto/cart"
	"time"

	"github.com/gofiber/fiber/v2"
)

func (c *ApiHandler) GetCart(ctx *fiber.Ctx) error {
	pagination := helper.PaginationNew(ctx)

	userID, ok := ctx.Locals("UserID").(string)
	if !ok {
		return helper.RespondHttpError(ctx, errors.New("user ID not found in context"))
	}

	context_, cancel := context.WithTimeout(ctx.UserContext(), time.Second*10)
	defer cancel()

	res, err := c.CartSvc.GetCartByUserID(context_, &cart.GetCartRequest{
		UserId: userID,
		Pagination: &cart.Pagination{
			Limit:  pagination.Limit,
			Offset: pagination.Offset,
		},
	})

	if err != nil {
		return helper.RespondHttpError(ctx, err)
	}

	res.Pagination.Page = int32(pagination.Page)

	storeItems := []*cart.StoreItems{}
	if res.StoreItems != nil {
		storeItems = res.StoreItems
	}

	return ctx.Status(200).JSON(fiber.Map{
		"items":       storeItems,
		"_pagination": res.Pagination,
	})
}

func (c *ApiHandler) AddItem(ctx *fiber.Ctx) error {
	request, err := helper.ParseAndValidateRequest(ctx, &constant.Products{})
	if err != nil {
		return helper.RespondHttpError(ctx, err)
	}

	userID, ok := ctx.Locals("UserID").(string)
	if !ok {
		return helper.RespondHttpError(ctx, errors.New("user ID not found in context"))
	}

	// godump.Dump(request)

	storeItems := []*cart.StoreItems{}
	for _, store := range request.Products {
		items := []*cart.CartItem{}
		for _, item := range store.Items {
			items = append(items, &cart.CartItem{
				ProductId:   item.ProductID,
				ProductName: item.ProductName,
				Quantity:    item.Quantity,
				Price:       item.Price,
				ImageUrl:    item.ImageURL,
			})
		}
		storeItems = append(storeItems, &cart.StoreItems{
			StoreID: store.StoreID,
			Items:   items,
		})
	}

	context_, cancel := context.WithTimeout(ctx.UserContext(), time.Second*10)
	defer cancel()

	res, err := c.CartSvc.AddItemToCart(context_, &cart.AddItemRequest{
		UserId:     userID,
		StoreItems: storeItems,
	})
	if err != nil {
		return helper.RespondHttpError(ctx, err)
	}

	return ctx.Status(200).JSON(fiber.Map{
		"message": res.Status,
	})
}

func (c *ApiHandler) RemoveItem(ctx *fiber.Ctx) error {
	request, err := helper.ParseAndValidateRequest(ctx, &constant.RemoveItemReq{}, helper.ParseOptions{Params: true})
	if err != nil {
		return helper.RespondHttpError(ctx, err)
	}

	userID, ok := ctx.Locals("UserID").(string)
	if !ok {
		return helper.RespondHttpError(ctx, errors.New("user ID not found in context"))
	}

	context_, cancel := context.WithTimeout(ctx.UserContext(), time.Second*10)
	defer cancel()

	res, err := c.CartSvc.RemoveItem(context_, &cart.RemoveFromCartRequest{
		UserId: userID,
		ItemId: int32(request.CartItemID),
	})
	if err != nil {
		return helper.RespondHttpError(ctx, err)
	}

	return ctx.Status(200).JSON(fiber.Map{
		"status": res.Status,
	})
}
