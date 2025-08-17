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
	version := "/api/v1"

	// Auth routes
	auth := app.Group(version + "/auth")
	auth.Post("/login", api.Login)
	auth.Post("/register", api.Register)

	// Cart routes
	cart := app.Group(version+"/cart", middleware.CheckRoles(constant.Admin))
	cart.Get("/", api.GetCart)
	cart.Post("/", api.AddItem)
	cart.Delete("/:cart_id/:cart_item_id", api.RemoveItem)

	// Store routes
	store := app.Group(version+"/store", middleware.CheckRoles(constant.Admin))
	store.Post("/create", api.CreateSrore)

	// Inventory routes
	inventory := app.Group(version+"/inventory", middleware.CheckRoles(constant.Admin))
	inventory.Post("/", api.AddInventory)
	inventory.Patch("/:product_id", api.UpdateInventory)
	inventory.Get("/", api.ListInventories)
	inventory.Delete("/:store_id/:inventory_id", api.RemoveInventory)

	// Catagory routes
	category := app.Group(version + "/catagory")
	category.Post("/", api.AddCategory)
	category.Patch("/", api.UpdateCategory)
	category.Get("/:store_id", api.GetCategory)

	// Payment routes
	payment := app.Group(version + "/payment")
	payment.Post("/webhook", api.StripeWebhook)
	payment.Post("/paid", middleware.CheckRoles(constant.Admin), api.Paid)

	// Order routes
	order := app.Group(version+"/order", middleware.CheckRoles(constant.Admin))
	order.Post("/", api.PlaceOrder)
}
