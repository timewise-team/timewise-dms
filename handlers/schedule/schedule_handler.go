package schedule

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/timewise-team/timewise-models/models"
	"gorm.io/gorm"
)

type ScheduleHandler struct {
	DB *gorm.DB
}

// GetSchedules godoc
// @Summary Get all schedules
// @Description Get all schedules
// @Tags schedule
// @Accept json
// @Produce json
// @Success 200 {array} models.TwSchedule
// @Router /dbms/v1/schedule [get]
func (h *ScheduleHandler) GetSchedules(c *fiber.Ctx) error {
	var schedules []models.TwSchedule
	if result := h.DB.Find(&schedules); result.Error != nil {
		// handle error
		return c.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}
	return c.JSON(schedules)
}

// GetScheduleById godoc
// @Summary Get schedule by ID
// @Description Get schedule by ID
// @Tags schedule
// @Accept json
// @Produce json
// @Param schedule_id path int true "Schedule ID"
// @Success 200 {object} models.TwSchedule
// @Router /dbms/v1/schedule/{schedule_id} [get]
func (h *ScheduleHandler) GetScheduleById(c *fiber.Ctx) error {
	var schedule models.TwSchedule
	scheduleId := c.Params("schedule_id")

	if err := h.DB.Where("id = ?", scheduleId).First(&schedule).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).SendString("Schedule not found")
		}
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	return c.JSON(schedule)
}
