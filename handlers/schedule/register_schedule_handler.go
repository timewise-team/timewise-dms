package schedule

import (
	"dbms/common"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func RegisterScheduleHandler(router fiber.Router, db *gorm.DB) {
	scheduleHandler := ScheduleHandler{
		DB: db,
	}
	common.RegisterHandler(router, db, func(handler common.Handler) {
		handler.Router.Get("/", scheduleHandler.GetSchedules)
		handler.Router.Get("/:schedule_id", scheduleHandler.GetScheduleById)
	})
}
