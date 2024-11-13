package reminder

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func RegisterReminderHandler(router fiber.Router, db *gorm.DB) {
	reminderHandler := ReminderHandler{
		DB: db,
	}
	router.Post("/", reminderHandler.CreateReminder)
	router.Get("/:reminder_id", reminderHandler.GetReminderById)
	router.Get("/schedule/:schedule_id", reminderHandler.GetRemindersByScheduleId)
	router.Put("/:reminder_id", reminderHandler.UpdateReminder)
	router.Delete("/:reminder_id", reminderHandler.DeleteReminder)
}
