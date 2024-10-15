package document

import (
	"github.com/gofiber/fiber/v2"
	"github.com/timewise-team/timewise-models/models"
)

// getDocumentsBySchedule godoc
// @Summary Get documents by schedule
// @Description Get documents by schedule
// @Tags document
// @Accept json
// @Produce json
// @Param schedule_id path string true "Schedule ID"
// @Success 200 {array} models.TwDocument
// @Router /dbms/v1/document/schedule/{schedule_id} [get]
func (h *DocumentHandler) getDocumentsBySchedule(c *fiber.Ctx) error {
	scheduleId := c.Params("schedule_id")
	if scheduleId == "" {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	var Documents []models.TwDocument
	if err := h.DB.
		Where("schedule_id = ?", scheduleId).
		Where("deleted_at IS NULL").
		Find(&Documents).Error; err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.JSON(Documents)
}
