package document

import (
	"github.com/gofiber/fiber/v2"
	"github.com/timewise-team/timewise-models/dtos/core_dtos/comment_dtos"
	"github.com/timewise-team/timewise-models/models"
	"log"
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

// getCommentsBySchedule godoc
// @Summary Get comments by schedule
// @Description Get comments by schedule
// @Tags comments
// @Accept json
// @Produce json
// @Param schedule_id path string true "Schedule ID"
// @Success 200 {array} models.TwComment
// @Router /dbms/v1/comment/schedule_id/{schedule_id} [get]
func (h *CommentHandler) getCommentsByScheduleID(c *fiber.Ctx) error {
	var scheduleComments []comment_dtos.TwCommentResponse
	scheduleId := c.Params("schedule_id")

	if scheduleId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Schedule ID không hợp lệ",
		})
	}

	// Perform the SQL query with multiple joins
	err := h.DB.Table("tw_comments AS c").
		Select(`
            c.id AS id,
            c.created_at,
            c.updated_at,
			c.schedule_id,
			c.workspace_user_id,
			c.commenter,
			c.content,
			c.is_deleted,
			wu.role,
			wu.status AS status_workspace_user,
			wu.is_verified,
			ue.id as user_id,
			ue.email,
			u.first_name,
			u.last_name,
			u.profile_picture
            
        `).
		Joins("JOIN tw_workspace_users AS wu ON wu.id =c.workspace_user_id").
		Joins("JOIN tw_user_emails AS ue ON wu.user_email_id = ue.id").
		Joins("JOIN tw_users AS u ON ue.user_id = u.id").
		Where("c.schedule_id = ?", scheduleId).
		Where("c.deleted_at IS NULL").
		Where("wu.deleted_at IS NULL").
		Where("ue.deleted_at IS NULL").
		Where("u.deleted_at IS NULL").
		Where("wu.is_active = true AND wu.is_verified = true AND wu.status = 'joined'").
		Scan(&scheduleComments).Error

	if err != nil {
		log.Println("Error querying schedule participants:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Không thể lấy danh sách participant",
		})
	}

	return c.JSON(scheduleComments)
}
