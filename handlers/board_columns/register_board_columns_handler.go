package board_columns

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type BoardColumnsHandler struct {
	Router fiber.Router
	DB     *gorm.DB
}

func RegisterBoardColumnsHandler(router fiber.Router, db *gorm.DB) {
	boardColumnsHandler := BoardColumnsHandler{
		Router: router,
		DB:     db,
	}

	// Register all endpoints here
	router.Get("/workspace/:workspace_id/board_columns", boardColumnsHandler.getBoardColumnsByWorkspace)
	router.Get("/workspace/:workspace_id/board_columns/:board_column_id", boardColumnsHandler.getBoardColumnById)
	router.Post("/board_columns", boardColumnsHandler.createBoardColumn)
	router.Put("/board_columns/:board_column_id", boardColumnsHandler.updateBoardColumn)
	router.Delete("/board_columns/:board_column_id", boardColumnsHandler.deleteBoardColumn)
	router.Get("/board_columns/:board_column_id/:field", boardColumnsHandler.getBoardColumnField)
	router.Put("/board_columns/:board_column_id/:field", boardColumnsHandler.updateBoardColumnField)

}
