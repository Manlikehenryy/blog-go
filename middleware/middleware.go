package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/manlikehenryy/blog-backend/util"
)

func IsAuthenticated(c *fiber.Ctx) error {
	cookie := c.Cookies("jwt")

	// Parse the JWT token from the cookie
	userId, err := util.ParseJwt(cookie)
	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	// Store the user ID in the request context
	c.Locals("userId", userId)

	return c.Next()
}