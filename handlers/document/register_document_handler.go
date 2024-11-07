package document

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type DocumentHandler struct {
	Router fiber.Router
	DB     *gorm.DB
}

func RegisterDocumentHandler(router fiber.Router, db *gorm.DB) {
	documentHandler := DocumentHandler{
		Router: router,
		DB:     db,
	}

	// Register all endpoints here
	router.Get("/schedule/:schedule_id", documentHandler.getDocumentsBySchedule)
	router.Get("/schedule_id/:schedule_id", documentHandler.getDocumentsByScheduleID)
	router.Post("/upload", documentHandler.createDocument)
}
