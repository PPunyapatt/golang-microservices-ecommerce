package handler

import (
	"context"
	"gateway/v1/internal/constant"
	"gateway/v1/internal/helper"
	"gateway/v1/proto/Inventory"
	"time"

	"github.com/gofiber/fiber/v2"
)

func (c *ApiHandler) AddInventory(ctx *fiber.Ctx) error {
	context_, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	var request *constant.Inventories
	if err := ctx.BodyParser(&request); err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	res, err := c.InventorySvc.AddInventory(context_, &Inventory.AddInvenRequest{
		Inventory: &Inventory.Inventory{
			StoreID:     request.StoreID,
			AddBy:       request.AddBy,
			Name:        request.Name,
			Description: request.Description,
			Price:       request.Price,
			Stock:       request.Stock,
			CategoryID:  request.CategoryID,
			ImageURL:    request.ImageURL,
		},
	})
	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{
			"error": "Failed to add inventory: " + err.Error(),
		})
	}

	return ctx.Status(200).JSON(res)
}

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
		Search:  request.Search,
		StoreID: request.StoreID,
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
