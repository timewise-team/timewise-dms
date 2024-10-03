package user

import (
	"errors"
	"github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	dtos "github.com/timewise-team/timewise-models/dtos/core_dtos"
	"github.com/timewise-team/timewise-models/dtos/core_dtos/user_login_dtos"
	"github.com/timewise-team/timewise-models/dtos/core_dtos/user_register_dtos"
	"github.com/timewise-team/timewise-models/models"
	"gorm.io/gorm"
	"time"
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
		var driverErr *mysql.MySQLError
		if errors.As(result.Error, &driverErr) && driverErr.Number == 1062 {
			return ctx.Status(fiber.StatusBadRequest).SendString("username already exists")
		}
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
// @Param user body dtos.UpdateUserRequest true "User object"
// @Success 200 {object} models.TwUser
// @Router /users/{id} [put]
func (h *UserHandler) updateUser(c *fiber.Ctx) error {
	var userDTO dtos.UpdateUserRequest
	if err := c.BodyParser(&userDTO); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	var user models.TwUser
	userId := c.Params("user_id")

	if err := h.DB.Where("id = ?", userId).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).SendString("User not found")
		}
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	// Update the fields if they are provided (not nil)
	if userDTO.FirstName != "" {
		user.FirstName = userDTO.FirstName
	}
	if userDTO.LastName != "" {
		user.LastName = userDTO.LastName
	}
	if userDTO.Email != "" {
		user.Email = userDTO.Email
	}
	if userDTO.ProfilePicture != "" {
		user.ProfilePicture = userDTO.ProfilePicture
	}
	if userDTO.Timezone != "" {
		user.Timezone = userDTO.Timezone
	}
	if userDTO.Locale != "" {
		user.Locale = userDTO.Locale
	}
	if userDTO.IsVerified {
		user.IsVerified = userDTO.IsVerified
	}
	if userDTO.IsActive {
		user.IsActive = userDTO.IsActive
	}
	if userDTO.NotificationSettings != "" {
		user.NotificationSettings = userDTO.NotificationSettings
	}
	if userDTO.CalendarSettings != "" {
		user.CalendarSettings = userDTO.CalendarSettings
	}
	if userDTO.Role != "" {
		user.Role = userDTO.Role
	}

	// Update the timestamp
	user.UpdatedAt = time.Now()

	if result := h.DB.Save(&user); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}

	return c.JSON(user)
}

func (h *UserHandler) deleteUser(ctx *fiber.Ctx) error {
	userId := ctx.Params("user_id")
	if result := h.DB.Delete(&models.TwUser{}, userId); result.Error != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}
	return ctx.SendString("User deleted successfully")
}

func (h *UserHandler) loginUser(ctx *fiber.Ctx) error {
	var user_login user_login_dtos.UserLoginRequest
	if err := ctx.BodyParser(&user_login); err != nil {
		return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	var user models.TwUser
	if err := h.DB.Where("username = ?", user_login.Username).First(&user).Error; err != nil {
		return ctx.Status(fiber.StatusUnauthorized).SendString("Invalid email or password")
	}

	return ctx.JSON(user)

}

func (h *UserHandler) getOrCreateUser(ctx *fiber.Ctx) error {
	// Parse the request body
	var req user_register_dto.GetOrCreateUserRequestDto
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	if !req.VerifiedEmail {
		return ctx.Status(fiber.StatusForbidden).SendString("User is not verified")
	}
	// Try to find the user in the database
	isNewUser := false
	var user models.TwUser
	if err := h.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// User not found, create a new one
			user = models.TwUser{
				Email:          req.Email,
				ProfilePicture: req.ProfilePicture,
				LastLoginAt:    time.Now(),
				Role:           "user",
				IsVerified:     req.VerifiedEmail,
				//GoogleId:       req.GoogleId,
				FirstName: req.GivenName,
				LastName:  req.FamilyName,
				Locale:    req.Locale,
				IsActive:  true,
			}

			if result := h.DB.Create(&user); result.Error != nil {
				return ctx.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
			}
			isNewUser = true
		} else {
			// Some other error occurred
			return ctx.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}
	}

	// Return the user
	resp := user_register_dto.GetOrCreateUserResponseDto{
		User:      user,
		IsNewUser: isNewUser,
	}
	return ctx.JSON(resp)
}
