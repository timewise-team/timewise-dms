package user_email

import (
	"errors"
	"github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"github.com/timewise-team/timewise-models/models"
	"gorm.io/gorm"
)

// @Summary Get all user emails
// @Description Get all user emails
// @Tags user_email
// @Accept json
// @Produce json
// @Success 200 {array} models.TwUserEmail
// @Router /dbms/v1/user_email [get]
func (h *UserEmailHandler) getUserEmails(c *fiber.Ctx) error {
	// Get user_id from query param
	userId := c.Query("user_id")
	if userId != "" {
		var userEmails []models.TwUserEmail
		if result := h.DB.Where("user_id = ?", userId).Find(&userEmails); result.Error != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
		}
		return c.JSON(userEmails)
	}
	var userEmails []models.TwUserEmail
	if result := h.DB.Find(&userEmails); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}

	return c.JSON(userEmails)
}

// @Summary Get user email by user ID
// @Description Get user email by user ID
// @Tags user_email
// @Accept json
// @Produce json
// @Param user_id path int true "User ID"
// @Success 200 {object} models.TwUserEmail
// @Router /dbms/v1/user_email/user/{user_id} [get]
func (h *UserEmailHandler) getUserEmailByUserId(c *fiber.Ctx) error {
	var userEmail models.TwUserEmail
	userId := c.Params("user_id")

	if err := h.DB.Where("user_id = ?", userId).First(&userEmail).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).SendString("Email not found")
		}
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	return c.JSON(userEmail)
}

// @Summary Create a new user email
// @Description Create a new user email
// @Tags user_email
// @Accept json
// @Produce json
// @Param email body models.TwUserEmail true "User Email"
// @Success 200 {object} models.TwUserEmail
// @Router /dbms/v1/user_email [post]
func (h *UserEmailHandler) createUserEmail(ctx *fiber.Ctx) error {
	userEmail := new(models.TwUserEmail)
	if err := ctx.BodyParser(userEmail); err != nil {
		return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	if result := h.DB.Create(&userEmail); result.Error != nil {
		var driverErr *mysql.MySQLError
		if errors.As(result.Error, &driverErr) && driverErr.Number == 1062 {
			return ctx.Status(fiber.StatusBadRequest).SendString("email already exists")
		}
		return ctx.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}
	// Lấy thông tin người dùng từ cơ sở dữ liệu dựa trên UserId
	var user models.TwUser
	if err := h.DB.Where("id = ?", userEmail.UserId).First(&user).Error; err != nil {
		return ctx.Status(fiber.StatusNotFound).SendString("User not found")
	}

	// Gán thông tin user vào userEmail để trả về kèm thông tin user
	userEmail.User = user

	return ctx.JSON(userEmail)
}

//func (h *UserEmailHandler) updateUserEmail(c *fiber.Ctx) error {
//	var userEmailDTO dtos.UpdateUserEmailRequest
//	if err := c.BodyParser(&userEmailDTO); err != nil {
//		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
//	}
//
//	var userEmail models.TwUserEmail
//	emailId := c.Params("email_id")
//
//	if err := h.DB.Where("id = ?", emailId).First(&userEmail).Error; err != nil {
//		if errors.Is(err, gorm.ErrRecordNotFound) {
//			return c.Status(fiber.StatusNotFound).SendString("Email not found")
//		}
//		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
//	}
//
//	// Update the fields if they are provided (not nil)
//	if userEmailDTO.Email != "" {
//		userEmail.Email = userEmailDTO.Email
//	}
//
//	if result := h.DB.Save(&userEmail); result.Error != nil {
//		return c.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
//	}
//
//	return c.JSON(userEmail)
//}

// @Summary Delete user email by ID
// @Description Delete user email by ID
// @Tags user_email
// @Accept json
// @Produce json
// @Param email_id path int true "Email ID"
// @Success 200 {string} string
// @Router /dbms/v1/user_email/{email_id} [delete]
func (h *UserEmailHandler) deleteUserEmail(ctx *fiber.Ctx) error {
	emailId := ctx.Params("email_id")
	if result := h.DB.Delete(&models.TwUserEmail{}, emailId); result.Error != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}
	return ctx.SendString("Email deleted successfully")
}
