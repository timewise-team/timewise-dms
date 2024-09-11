package user

import (
	"github.com/gofiber/fiber/v2"
	"github.com/timewise-team/timewise-models/models"
)

// GET /users
// getUsers godoc
// @Summary Get all users
// @Description Get all users
// @Tags user
// @Accept json
// @Produce json
// @Success 200 {array} models.TwUser
// @Router /users [get]
func (h *UserHandler) getUsers(c *fiber.Ctx) error {
	var users []models.TwUser
	if result := h.DB.Find(&users); result.Error != nil {
		// handle error
		return c.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}

	return c.JSON(users)
}
