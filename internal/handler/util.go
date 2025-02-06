package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// getCurrentUserID extracts the authenticated user's ID from the Fiber context
func getCurrentUserID(c *fiber.Ctx) (uuid.UUID, error) {
	userID, ok := c.Locals("userID").(string)
	if !ok {
		return uuid.Nil, fiber.NewError(fiber.StatusUnauthorized, "User ID not found in context")
	}

	id, err := uuid.Parse(userID)
	if err != nil {
		return uuid.Nil, fiber.NewError(fiber.StatusInternalServerError, "Invalid user ID format")
	}

	return id, nil
}
