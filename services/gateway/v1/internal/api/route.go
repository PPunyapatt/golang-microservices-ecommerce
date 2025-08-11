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

	// Auth routes
	auth := app.Group("/api/v1/auth")
	auth.Post("/login", api.Login)
	auth.Post("/register", api.Register)

	// Cart routes
	cart := app.Group("/api/v1/cart", middleware.CheckRoles(constant.Admin))
	cart.Get("/", api.GetCart)
	cart.Post("/", api.AddItem)
	cart.Delete("/:cart_id/:cart_item_id", api.RemoveItem)

	// Store routes
	store := app.Group("/api/v1/store", middleware.CheckRoles(constant.Admin))
	store.Post("/create", api.CreateSrore)

	// Inventory routes
	inventory := app.Group("/api/v1/inventory", middleware.CheckRoles(constant.Admin))
	inventory.Post("/", api.AddInventory)
	inventory.Patch("/", api.UpdateInventory)
	inventory.Get("/", api.ListInventories)
	inventory.Delete("/:store_id/:inventory_id", api.RemoveInventory)

	// Catagory routes
	category := app.Group("/api/v1/catagory")
	category.Post("/", api.AddCategory)
	category.Patch("/", api.UpdateCategory)
	category.Get("/:store_id", api.GetCategory)

	// Payment routes
	payment := app.Group("/api/v1/payment")
	payment.Post("/webhook", api.StripeWebhook)
	payment.Post("/paid", middleware.CheckRoles(constant.Admin), api.Paid)
}
