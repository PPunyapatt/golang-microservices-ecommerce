package handler

import (
	"context"
	"errors"
	"gateway/v1/internal/constant"
	"gateway/v1/internal/helper"
	"gateway/v1/proto/Inventory"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/goforj/godump"
)

// ------------ Inventory --------------

func (c *ApiHandler) AddInventory(ctx *fiber.Ctx) error {
	request, err := helper.ParseAndValidateRequest(ctx, &constant.Inventories{})
	if err != nil {
		return helper.RespondHttpError(ctx, err)
	}

	godump.Dump(request)

	context_, cancel := context.WithTimeout(ctx.UserContext(), time.Second*10)
	defer cancel()

	userID, ok := ctx.Locals("UserID").(string)
	if !ok {
		return helper.RespondHttpError(ctx, errors.New("user ID not found in context"))
	}

	res, err := c.InventorySvc.AddInventory(context_, &Inventory.AddInvenRequest{
		Inventory: &Inventory.Inventory{
			StoreID:     request.StoreID,
			AddBy:       &userID,
			Name:        request.Name,
			Description: request.Description,
			Price:       request.Price,
			Stock:       request.AvailableStock,
			CategoryID:  request.CategoryID,
			ImageURL:    request.ImageURL,
		},
	})
	if err != nil {
		return helper.RespondHttpError(ctx, err)
	}

	return ctx.Status(200).JSON(res)
}

func (c *ApiHandler) UpdateInventory(ctx *fiber.Ctx) error {
	request, err := helper.ParseAndValidateRequest(ctx, &constant.Inventories{}, helper.ParseOptions{Params: true, Body: true})
	if err != nil {
		return helper.RespondHttpError(ctx, err)
	}
	// godump.Dump(request)

	context_, cancel := context.WithTimeout(ctx.UserContext(), time.Second*10)
	defer cancel()

	res, err := c.InventorySvc.UpdateInventory(context_, &Inventory.UpdateInvenRequest{
		Inventory: &Inventory.Inventory{
			ID:          request.ID,
			StoreID:     request.StoreID,
			AddBy:       request.AddBy,
			Name:        request.Name,
			Description: request.Description,
			Price:       request.Price,
			Stock:       request.AvailableStock,
			CategoryID:  request.CategoryID,
			ImageURL:    request.ImageURL,
		},
	})
	if err != nil {
		return helper.RespondHttpError(ctx, err)
	}

	return ctx.Status(200).JSON(fiber.Map{
		"message": res.Status,
	})
}

func (c *ApiHandler) ListInventories(ctx *fiber.Ctx) error {
	request, err := helper.ParseAndValidateRequest(ctx, &constant.GetCategoryReq{}, helper.ParseOptions{Query: true})
	if err != nil {
		return helper.RespondHttpError(ctx, err)
	}
	pagination := helper.PaginationNew(ctx)

	context_, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	res, err := c.InventorySvc.ListInventories(context_, &Inventory.ListInvetoriesRequest{
		Fields: &Inventory.Search{
			Query:      request.Query,
			CategoryID: request.CategoryID,
			StoreID:    request.StoreID,
		},
		Pagination: &Inventory.Pagination{
			Limit:  pagination.Limit,
			Offset: pagination.Offset,
		},
	})

	if err != nil {
		return helper.RespondHttpError(ctx, err)
	}

	res.Pagination.Page = int32(pagination.Page)

	return ctx.Status(200).JSON(fiber.Map{
		"data":        res.Inventory,
		"_pagination": res.Pagination,
	})
}

func (c *ApiHandler) RemoveInventory(ctx *fiber.Ctx) error {
	request, err := helper.ParseAndValidateRequest(ctx, &constant.RemoveInventoryReq{}, helper.ParseOptions{Params: true})
	if err != nil {
		return helper.RespondHttpError(ctx, err)
	}

	userID, ok := ctx.Locals("UserID").(string)
	if !ok {
		return helper.RespondHttpError(ctx, errors.New("user ID not found in context"))
	}

	context_, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	res, err := c.InventorySvc.RemoveInventory(context_, &Inventory.RemoveInvenRequest{
		InvetoriesID: int32(request.InventoryID),
		StoreID:      int32(request.StoreID),
		UserID:       userID,
	})
	if err != nil {
		return helper.RespondHttpError(ctx, err)
	}

	return ctx.Status(200).JSON(fiber.Map{
		"status": res.Status,
	})
}

// ------------ Category --------------

func (c *ApiHandler) AddCategory(ctx *fiber.Ctx) error {
	context_, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	catagoryName := ctx.Query("name")

	res, err := c.InventorySvc.AddCategory(context_, &Inventory.AddCategoryRequest{
		Name: catagoryName,
	})
	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{
			"error": "Failed to add category: " + err.Error(),
		})
	}

	return ctx.Status(200).JSON(fiber.Map{
		"message": res.Status,
	})
}

func (c *ApiHandler) UpdateCategory(ctx *fiber.Ctx) error {
	context_, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	var request *constant.Category
	if err := ctx.BodyParser(&request); err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	res, err := c.InventorySvc.UpdateCategory(context_, &Inventory.UpdateCategoryRequest{
		Name:       request.Name,
		StoreID:    request.StoreID,
		CategoryID: request.ID,
	})
	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{
			"error": "Failed to update category: " + err.Error(),
		})
	}

	return ctx.Status(200).JSON(fiber.Map{
		"message": res.Status,
	})
}

func (c *ApiHandler) GetCategory(ctx *fiber.Ctx) error {
	pagination := helper.PaginationNew(ctx)
	request, err := helper.ParseAndValidateRequest(ctx, &constant.GetCategoryReq{}, helper.ParseOptions{Params: true})
	if err != nil {
		return helper.RespondHttpError(ctx, err)
	}

	context_, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	res, err := c.InventorySvc.ListCategories(context_, &Inventory.ListCategoriesRequest{
		Pagination: &Inventory.Pagination{
			Limit:  pagination.Limit,
			Offset: pagination.Offset,
		},
		Fields: &Inventory.Search{
			Query:      request.Query,
			CategoryID: request.CategoryID,
			StoreID:    request.StoreID,
		},
	})
	if err != nil {
		helper.RespondHttpError(ctx, err)
	}

	if res.Catagories == nil {
		res.Catagories = []*Inventory.Category{}
	}

	res.Pagination.Page = int32(pagination.Page)

	return ctx.Status(200).JSON(fiber.Map{
		"_pagination": res.Pagination,
		"Data":        res.Catagories,
	})
}
