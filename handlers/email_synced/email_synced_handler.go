package email_synced

import (
	"github.com/gofiber/fiber/v2"
	"github.com/timewise-team/timewise-models/models"
	"gorm.io/gorm"
	"net/url"
)

type EmailSyncedHandler struct {
	DB *gorm.DB
}

// GetAllEmailSynced godoc
// @Summary Get all emails synced with an email
// @Description Get all emails synced with an email
// @Tags Email Synced
// @Accept json
// @Produce json
// @Param email path string true "Email"
// @Success 200 {object} models.TwUserEmail "List of emails synced with the email"
// @Failure 400 {object} fiber.Map "Invalid email"
// @Failure 400 {object} fiber.Map "Email is not synced any other emails"
// @Router /dbms/v1/email_synced/{email} [get]
func (e *EmailSyncedHandler) GetAllEmailSynced(c *fiber.Ctx) error {
	encodedEmail := c.Params("email")
	if encodedEmail == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Email is required",
		})
	}
	email, err := url.QueryUnescape(encodedEmail) // Decode '%40' into '@'
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid email format",
		})
	}
	// check if email existed
	if !e.checkEmailExisted(email) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Email is not existed",
		})
	}
	var emails []models.TwUserEmail
	result := e.DB.Where("email = ?", email).Find(&emails)

	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Account is not synced with any other emails",
		})
	}
	return c.JSON(emails)
}

// function to check if email existed
func (e *EmailSyncedHandler) checkEmailExisted(email string) bool {
	if e.DB.Where("email = ?", email).Find(&models.TwUser{}).RowsAffected == 0 {
		return false
	}
	return true
}
