package api

import (
	"gateway/v1/internal/api/handler"

	"github.com/gofiber/fiber/v2"
)

func Route(
	app *fiber.App,
	api *handler.ApiHandler,
) {

	// Initialize the CartHandler with the cart service client

	// Cart routes
	cart := app.Group("/api/v1/cart")
	cart.Post("/add", api.AddItem)
	cart.Post("/remove", api.RemoveItem)

	// Auth routes
	auth := app.Group("/api/v1/auth")
	auth.Post("/login", api.Login)
}
