package workspace_log

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/timewise-team/timewise-models/models"
	"gorm.io/gorm"
)

type WorkspaceLog struct {
	Router fiber.Router
	DB     *gorm.DB
}

func RegisterWorkspaceLogHandler(router fiber.Router, db *gorm.DB) {
	workspaceLogHandler := WorkspaceLog{
		Router: router,
		DB:     db,
	}

	// Register all endpoints here
	router.Get("/", workspaceLogHandler.getWorkspaceLog)
	router.Get("/:workspace_log_id", workspaceLogHandler.getWorkspaceLogById)
	router.Delete("/:workspace_log_id", workspaceLogHandler.removeWorkspaceLogById)
}

func (h *WorkspaceLog) getWorkspaceLog(c *fiber.Ctx) error {
	var workspaceLogs []models.TwWorkspaceLog
	if result := h.DB.Find(&workspaceLogs); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}

	return c.JSON(workspaceLogs)
}

func (h *WorkspaceLog) getWorkspaceLogById(c *fiber.Ctx) error {
	var workspaceLog models.TwWorkspaceLog
	workspaceLogId := c.Params("workspace_log_id")

	if err := h.DB.Where("id = ?", workspaceLogId).First(&workspaceLog).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).SendString("WorkspaceLog not found")
		}
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	return c.JSON(workspaceLog)
}
func (h *WorkspaceLog) removeWorkspaceLogById(c *fiber.Ctx) error {
	workspaceLogId := c.Params("workspace_log_id")
	if result := h.DB.Delete(&models.TwWorkspaceLog{}, workspaceLogId); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}
	return c.SendStatus(fiber.StatusNoContent)

}
