package user

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	dtos "github.com/timewise-team/timewise-models/dtos/core_dtos"
	"github.com/timewise-team/timewise-models/models"
	"gorm.io/gorm"
	"reflect"
)

// GET /users/{id}
// getUserById godoc
// @Summary Get user by ID
// @Description Get user by ID
// @Tags user
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} models.TwUser
// @Router /users/{id} [get]
func (h *UserHandler) getUserById(c *fiber.Ctx) error {
	var user models.TwUser
	userId := c.Params("user_id")

	if err := h.DB.Where("id = ?", userId).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).SendString("User not found")
		}
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	return c.JSON(user)
}

// POST /users
// createUser godoc
// @Summary Create a new user
// @Description Create a new user
// @Tags user
// @Accept json
// @Produce json
// @Param user body models.TwUser true "User object"
// @Success 200 {object} models.TwUser
// @Router /users [post]
func (h *UserHandler) createUser(ctx *fiber.Ctx) error {
	user := new(models.TwUser)
	if err := ctx.BodyParser(user); err != nil {
		return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	if result := h.DB.Create(&user); result.Error != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}

	return ctx.JSON(user)
}

// PUT /users/{id}
// updateUser godoc
// @Summary Update user by ID
// @Description Update user by ID
// @Tags user
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param user body models.TwUser true "User object"
// @Success 200 {object} models.TwUser
// @Router /users/{id} [put]
func (h *UserHandler) updateUser(c *fiber.Ctx) error {
	// Parse the body data into a UpdateUserRequest DTO
	var req dtos.UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	updatedFields := req.UpdatedFields

	// Create a map for the Updates method
	updates := make(map[string]interface{})
	for _, field := range updatedFields {
		// Use reflection to get the field value from the TwUser object
		r := reflect.ValueOf(&req.User).Elem()
		f := r.FieldByName(field)
		if f.IsValid() {
			fieldValue := f.Interface()
			updates[field] = fieldValue
		} else {
			return c.Status(fiber.StatusBadRequest).SendString("Invalid field name: " + field)
		}
	}

	// Perform the update operation
	if result := h.DB.Model(&req.User).Where("id = ?", req.User.ID).Updates(updates); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}

	return c.JSON(req.User)
}
