package notification

import (
	"dbms/common"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func RegisterNotificationHandler(router fiber.Router, db *gorm.DB) {
	notification := NotificationHandler{
		DB: db,
	}
	common.RegisterHandler(router, db, func(handler common.Handler) {
		handler.Router.Post("/", notification.CreateNotification)
	})
}
