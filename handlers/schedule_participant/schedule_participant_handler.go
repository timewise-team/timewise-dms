package schedule_participant

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/timewise-team/timewise-models/models"
	"gorm.io/gorm"
)

func (h *ScheduleParticipantHandler) getScheduleParticipants(c *fiber.Ctx) error {
	var scheduleParticipants []models.TwScheduleParticipant
	if result := h.DB.Find(&scheduleParticipants); result.Error != nil {
		// handle error
		return c.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}

	return c.JSON(scheduleParticipants)
}

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
