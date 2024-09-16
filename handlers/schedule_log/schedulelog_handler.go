package schedule_log

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/timewise-team/timewise-models/models"
	"gorm.io/gorm"
)

func (h *ScheduleLogHandler) getScheduleLogs(c *fiber.Ctx) error {
	var scheduleLogs []models.TwScheduleLog
	if result := h.DB.Find(&scheduleLogs); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}

	return c.JSON(scheduleLogs)
}

func (h *ScheduleLogHandler) getScheduleLogById(c *fiber.Ctx) error {
	var scheduleLog models.TwScheduleLog
	scheduleLogId := c.Params("id")

	if err := h.DB.Where("id = ?", scheduleLogId).First(&scheduleLog).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).SendString("ScheduleLog not found")
		}
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	return c.JSON(scheduleLog)
}

func (h *ScheduleLogHandler) updateScheduleLog(c *fiber.Ctx) error {
	var scheduleLog models.TwScheduleLog

	if result := h.DB.First(&scheduleLog, c.Params("id")); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}

	if err := c.BodyParser(&scheduleLog); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	if result := h.DB.Save(&scheduleLog); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}
	return c.JSON(scheduleLog)
}

func (h *ScheduleLogHandler) deleteScheduleLog(c *fiber.Ctx) error {
	var scheduleLog models.TwScheduleLog
	if result := h.DB.First(&scheduleLog, c.Params("id")); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}
	if result := h.DB.Delete(&scheduleLog); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}
	return c.JSON(fiber.Map{
		"status": "success",
	})
}

func (h *ScheduleLogHandler) createScheduleLog(c *fiber.Ctx) error {
	var scheduleLog models.TwScheduleLog
	if err := c.BodyParser(&scheduleLog); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	if result := h.DB.Create(&scheduleLog); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}
	return c.JSON(scheduleLog)
}
