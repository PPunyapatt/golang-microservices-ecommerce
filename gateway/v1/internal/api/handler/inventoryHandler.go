package handler

import (
	"context"
	"gateway/v1/internal/constant"
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
			CatagoryID:  request.CategoryID,
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
		CatagoryID: request.ID,
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
