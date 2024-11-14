package notification_setting

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type NotificationSettingHandler struct {
	DB *gorm.DB
}

func (h NotificationSettingHandler) GetNotificationSettings(ctx *fiber.Ctx) error {
	return nil
}

func (h NotificationSettingHandler) GetNotificationSettingByUserId(ctx *fiber.Ctx) error {
	return nil

}

func (h NotificationSettingHandler) CreateNotificationSetting(ctx *fiber.Ctx) error {
	return nil

}

func (h NotificationSettingHandler) UpdateNotificationSetting(ctx *fiber.Ctx) error {
	return nil

}
