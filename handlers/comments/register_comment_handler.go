package document

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type CommentHandler struct {
	Router fiber.Router
	DB     *gorm.DB
}

func RegisterCommentsHandler(router fiber.Router, db *gorm.DB) {
	commentHandler := CommentHandler{
		Router: router,
		DB:     db,
	}

	// Register all endpoints here
	router.Get("/schedule/:schedule_id", commentHandler.getCommentsBySchedule)
	router.Get("/schedule_id/:schedule_id", commentHandler.getCommentsByScheduleID)

}
