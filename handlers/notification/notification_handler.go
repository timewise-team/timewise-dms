package notification

import (
	"errors"
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
// @Router /dbms/v1/notification [post]
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

// GetUnsentNotifications godoc
// @Summary Get unsent notifications
// @Description Get unsent notifications
// @Tags notification
// @Accept json
// @Produce json
// @Success 200 {array} models.TwNotifications
// @Router /dbms/v1/notification [get]
func (h *NotificationHandler) GetUnsentNotifications(ctx *fiber.Ctx) error {
	var notifications []models.TwNotifications
	if err := h.DB.Where("is_sent = ?", false).Preload("UserEmail").Find(&notifications).Error; err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return ctx.JSON(notifications)
}

// updateNotificationToSent godoc
// @Summary Update notification to sent
// @Description Update notification to sent
// @Tags notification
// @Accept json
// @Produce json
// @Param notification_id path string true "Notification ID"
// @Success 200 {object} fiber.Map
// @Router /dbms/v1/notification/{notification_id} [put]
func (h *NotificationHandler) updateNotificationToSent(ctx *fiber.Ctx) error {
	notificationID := ctx.Params("notification_id")
	if notificationID == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Notification ID is required",
		})
	}

	var notification models.TwNotifications
	if err := h.DB.First(&notification, notificationID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Notification not found",
			})
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	notification.IsSent = true
	if err := h.DB.Save(&notification).Error; err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.JSON(fiber.Map{
		"message": "Notification updated to sent",
	})
}
