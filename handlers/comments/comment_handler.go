package document

import (
	"github.com/gofiber/fiber/v2"
	"github.com/timewise-team/timewise-models/models"
)

// getCommentsBySchedule godoc
// @Summary Get comments by schedule
// @Description Get comments by schedule
// @Tags comments
// @Accept json
// @Produce json
// @Param schedule_id path string true "Schedule ID"
// @Success 200 {array} models.TwComment
// @Router /dbms/v1/comment/schedule/{schedule_id} [get]
func (h *CommentHandler) getCommentsBySchedule(c *fiber.Ctx) error {
	scheduleId := c.Params("schedule_id")
	if scheduleId == "" {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	var Comments []models.TwComment
	if err := h.DB.
		Where("schedule_id = ?", scheduleId).
		Where("deleted_at IS NULL").
		Find(&Comments).Error; err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.JSON(Comments)
}
