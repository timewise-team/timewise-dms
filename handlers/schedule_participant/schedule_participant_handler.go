package schedule_participant

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/timewise-team/timewise-models/models"
	"gorm.io/gorm"
)

// getScheduleParticipants godoc
// @Summary Get all schedule participants
// @Description Get all schedule participants
// @Tags schedule_participant
// @Accept json
// @Produce json
// @Success 200 {array} models.TwScheduleParticipant
// @Router /dbms/v1/schedule_participant [get]
func (h *ScheduleParticipantHandler) getScheduleParticipants(c *fiber.Ctx) error {
	var scheduleParticipants []models.TwScheduleParticipant
	if result := h.DB.Find(&scheduleParticipants); result.Error != nil {
		// handle error
		return c.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}

	return c.JSON(scheduleParticipants)
}

// @Summary Get schedule participant by ID
// @Description Get schedule participant by ID
// @Tags schedule_participant
// @Accept json
// @Produce json
// @Param id path int true "Schedule Participant ID"
// @Success 200 {object} models.TwScheduleParticipant
// @Router /dbms/v1/schedule_participant/{id} [get]
func (h *ScheduleParticipantHandler) getScheduleParticipantById(c *fiber.Ctx) error {
	var scheduleParticipant models.TwScheduleParticipant
	scheduleParticipantId := c.Params("id")

	if err := h.DB.Where("id = ?", scheduleParticipantId).First(&scheduleParticipant).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).SendString("ScheduleParticipant not found")
		}
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	return c.JSON(scheduleParticipant)
}

// updateScheduleParticipant godoc
// @Summary Update schedule participant
// @Description Update schedule participant
// @Tags schedule_participant
// @Accept json
// @Produce json
// @Param id path int true "Schedule Participant ID"
// @Param schedule_participant body models.TwScheduleParticipant true "Schedule participant object"
// @Success 200 {object} models.TwScheduleParticipant
// @Router /dbms/v1/schedule_participant/{id} [put]
func (h *ScheduleParticipantHandler) updateScheduleParticipant(c *fiber.Ctx) error {
	var scheduleParticipants models.TwScheduleParticipant
	if result := h.DB.First(&scheduleParticipants, c.Params("id")); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}
	if err := c.BodyParser(&scheduleParticipants); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	if result := h.DB.Save(&scheduleParticipants); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}
	return c.JSON(scheduleParticipants)

}

// deleteScheduleParticipant godoc
// @Summary Delete schedule participant
// @Description Delete schedule participant
// @Tags schedule_participant
// @Accept json
// @Produce json
// @Param id path int true "Schedule Participant ID"
// @Success 200 {object} models.TwScheduleParticipant
// @Router /dbms/v1/schedule_participant/{id} [delete]
func (h *ScheduleParticipantHandler) deleteScheduleParticipant(c *fiber.Ctx) error {
	var scheduleParticipants models.TwScheduleParticipant
	if result := h.DB.First(&scheduleParticipants, c.Params("id")); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}
	if result := h.DB.Delete(&scheduleParticipants); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}
	return c.JSON(fiber.Map{
		"status": "deleted",
	})
}

// createScheduleParticipant godoc
// @Summary Create schedule participant
// @Description Create schedule participant
// @Tags schedule_participant
// @Accept json
// @Produce json
// @Param schedule_participant body models.TwScheduleParticipant true "Schedule participant object"
// @Success 200 {object} models.TwScheduleParticipant
// @Router /dbms/v1/schedule_participant [post]
func (h *ScheduleParticipantHandler) createScheduleParticipant(c *fiber.Ctx) error {
	var scheduleParticipants models.TwScheduleParticipant
	if err := c.BodyParser(&scheduleParticipants); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	if result := h.DB.Create(&scheduleParticipants); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}
	return c.JSON(scheduleParticipants)
}
