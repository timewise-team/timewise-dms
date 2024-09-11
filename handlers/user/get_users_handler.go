package user

import (
	"github.com/gofiber/fiber/v2"
	"github.com/timewise-team/timewise-models/models"
)

// GET /users
func (h *UserHandler) getUsers(c *fiber.Ctx) error {
	var users []models.TwUser
	if result := h.DB.Find(&users); result.Error != nil {
		// handle error
		return c.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}

	return c.JSON(users)
}
