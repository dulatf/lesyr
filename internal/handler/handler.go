package handler

import (
	"github.com/gofiber/fiber/v2"
)

// ErrorHandler handles all errors globally
func ErrorHandler(ctx *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError

	// Check if it's a fiber error
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}

	return ctx.Status(code).JSON(fiber.Map{
		"error": err.Error(),
	})
}

// HealthCheck handles the health check endpoint
func HealthCheck(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status":  "OK",
		"message": "Server is healthy",
	})
}

// AuthMiddleware placeholder - we'll implement this properly later
func AuthMiddleware(c *fiber.Ctx) error {
	// TODO: Implement proper JWT validation
	return c.Next()
}

func CreateFeed(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "Create feed - Not implemented"})
}

func GetFeed(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "Get feed - Not implemented"})
}

func DeleteFeed(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "Delete feed - Not implemented"})
}

// Auth handlers - we'll implement these properly later
func Register(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "Register - Not implemented"})
}

func Login(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "Login - Not implemented"})
}
