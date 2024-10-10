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
		handler.Router.Get("/schedules/filter", scheduleHandler.FilterSchedules)
		//handler.Router.Get("/user/:user_id", scheduleHandler.GetSchedulesByUserId)
		handler.Router.Post("/", scheduleHandler.CreateSchedule)
		handler.Router.Put("/:schedule_id", scheduleHandler.UpdateSchedule)
		handler.Router.Delete("/:schedule_id", scheduleHandler.DeleteSchedule)
		router.Get("/workspace/:workspace_id/board_columns/:board_column_id/schedules", scheduleHandler.getSchedulesByBoardColumn)
		router.Get("/workspace/:workspace_id/schedules", scheduleHandler.GetSchedulesByWorkspace)
	})
}
