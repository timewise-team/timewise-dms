package schedule_participant

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type ScheduleParticipantHandler struct {
	Router fiber.Router
	DB     *gorm.DB
}

func RegisterScheduleParticipantHandler(router fiber.Router, db *gorm.DB) {
	scheduleParticipantHandeler := ScheduleParticipantHandler{
		Router: router,
		DB:     db,
	}

	// Register all endpoints here
	router.Get("/", scheduleParticipantHandeler.getScheduleParticipants)
	router.Get("/:id", scheduleParticipantHandeler.getScheduleParticipantById)
	router.Post("/", scheduleParticipantHandeler.createScheduleParticipant)
	router.Put("/:id", scheduleParticipantHandeler.updateScheduleParticipant)
	router.Delete("/:id", scheduleParticipantHandeler.deleteScheduleParticipant)
}
