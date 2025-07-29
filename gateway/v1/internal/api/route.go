package api

import (
	"gateway/v1/internal/api/handler"
	"gateway/v1/internal/constant"
	"gateway/v1/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

func Route(
	app *fiber.App,
	api *handler.ApiHandler,
) {

	// Initialize the CartHandler with the cart service client

	// Cart routes
	cart := app.Group("/api/v1/cart", middleware.CheckRoles(constant.Admin))
	cart.Post("/add", api.AddItem)
	cart.Post("/remove", api.RemoveItem)

	// Auth routes
	auth := app.Group("/api/v1/auth")
	auth.Post("/login", api.Login)
	auth.Post("/register", api.Register)

	// Store routes
	store := app.Group("/api/v1/store", middleware.CheckRoles(constant.Admin))
	store.Post("/create", api.CreateSrore)

	// Inventory routes
	inventory := app.Group("/api/v1/inventory", middleware.CheckRoles(constant.Admin))
	inventory.Post("/", api.AddInventory)

	// Catagory routes
	catagory := app.Group("/api/v1/catagory")
	catagory.Post("/", api.AddCategory)
	catagory.Patch("/", api.UpdateCategory)
}
