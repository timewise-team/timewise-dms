package document

import (
	"github.com/gofiber/fiber/v2"
	"github.com/timewise-team/timewise-models/dtos/core_dtos/document_dtos"
	"log"
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
	var documents []document_dtos.TwDocumentResponse
	err := h.DB.Table("tw_documents AS d").
		Select(`
            d.id AS id,
            d.file_name,
            d.file_path,
			d.file_size,
			d.file_type,
			d.is_deleted,
			d.schedule_id,
			d.uploaded_by,
			d.created_at,
			d.updated_at,
			d.uploaded_at,
			d.download_url,
			wu.role,
			wu.status AS status_workspace_user,
			wu.is_verified,
			ue.id as user_id,
			ue.email,
			u.first_name,
			u.last_name,
			u.profile_picture
            
        `).
		Joins("JOIN tw_workspace_users AS wu ON wu.id =d.uploaded_by").
		Joins("JOIN tw_user_emails AS ue ON wu.user_email_id = ue.id").
		Joins("JOIN tw_users AS u ON ue.user_id = u.id").
		Where("d.schedule_id = ?", scheduleId).
		Where("d.deleted_at IS NULL").
		Where("wu.deleted_at IS NULL").
		Where("ue.deleted_at IS NULL").
		Where("u.deleted_at IS NULL").
		Where("wu.is_active = true AND wu.is_verified = true AND wu.status = 'joined'").
		Scan(&documents).Error
	if err != nil {
		log.Println("Error querying document:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Không thể lấy danh sách document",
		})
	}

	return c.JSON(documents)
}
