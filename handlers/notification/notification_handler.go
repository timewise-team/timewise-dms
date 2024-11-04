package notification

import (
	"github.com/gofiber/fiber/v2"
	"github.com/timewise-team/timewise-models/dtos/core_dtos"
	"github.com/timewise-team/timewise-models/models"
	"gorm.io/gorm"
)

type NotificationHandler struct {
	DB *gorm.DB
}

// CreateNotification godoc
// @Summary Create a new notification
// @Description Create a new notification
// @Tags notification
// @Accept json
// @Produce json
// @Param createNotificationRequest body core_dtos.PushNotificationDto true "Create notification request"
// @Success 200 {object} core_dtos.PushNotificationDto
// @Router /api/v1/notification [post]
func (h *NotificationHandler) CreateNotification(c *fiber.Ctx) error {
	// Get data from request
	var request core_dtos.PushNotificationDto
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	notification := models.TwNotifications{
		UserEmailId:     request.UserEmailId,
		Type:            request.Type,
		Message:         request.Message,
		RelatedItemId:   request.RelatedItemId,
		RelatedItemType: request.RelatedItemType,
		ExtraData:       request.ExtraData,
	}
	if result := h.DB.Create(&notification); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": result.Error.Error()})
	}
	// Insert data into database
	return nil
}
