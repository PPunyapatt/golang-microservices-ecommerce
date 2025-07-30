package helper

import (
	"github.com/gofiber/fiber/v2"
)

func RespondHttpError(ctx *fiber.Ctx, err error) error {
	// Respond with a generic error message
	return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"error": err.Error(),
	})
}
