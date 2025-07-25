package middleware

import (
	"gateway/v1/internal/constant"
	"gateway/v1/internal/helper"

	"github.com/gofiber/fiber/v2"
)

func CheckRoles(allowedRoles ...constant.Role) fiber.Handler {
	return func(context *fiber.Ctx) error {
		tokenString := context.Get("Authorization")
		claims, err := helper.VerifyToken(tokenString)
		if err != nil {
			return context.Status(503).JSON(fiber.Map{
				"error": "Failed to verify token" + err.Error(),
			})
		}
		context.Locals("UserID", claims.ID)

		for _, allowedRole := range allowedRoles {
			for _, role := range claims.Roles {
				if constant.Role(role) == allowedRole {
					return context.Next()
				}
			}
		}

		return context.Status(503).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
}
