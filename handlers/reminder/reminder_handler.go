package reminder

import (
	"github.com/gofiber/fiber/v2"
	"github.com/timewise-team/timewise-models/models"
	"gorm.io/gorm"
)

type ReminderHandler struct {
	DB *gorm.DB
}

// createReminder godoc
// @Summary Create a reminder
// @Description Create a reminder
// @Tags reminder
// @Accept json
// @Produce json
// @Param reminder body models.TwReminder true "Reminder"
// @Success 200 {object} models.TwReminder
// @Router /api/v1/reminder [post]
func (h ReminderHandler) CreateReminder(ctx *fiber.Ctx) error {
	var reminder models.TwReminder
	if err := ctx.BodyParser(&reminder); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	if result := h.DB.Create(&reminder); result.Error != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}
	return ctx.JSON(reminder)
}

// getReminderById godoc
// @Summary Get reminder by ID
// @Description Get reminder by ID
// @Tags reminder
// @Accept json
// @Produce json
// @Param reminder_id path string true "Reminder ID"
// @Success 200 {object} models.TwReminder
// @Router /api/v1/reminder/{reminder_id} [get]
func (h ReminderHandler) GetReminderById(ctx *fiber.Ctx) error {
	id := ctx.Params("reminder_id")
	var reminder models.TwReminder
	if result := h.DB.
		Where("id = ?", id).
		Where("deleted_at IS NULL").
		First(&reminder); result.Error != nil {
		return ctx.Status(fiber.StatusNotFound).SendString(result.Error.Error())
	}
	return ctx.JSON(reminder)
}

// getRemindersByScheduleId godoc
// @Summary Get reminders by schedule ID
// @Description Get reminders by schedule ID
// @Tags reminder
// @Accept json
// @Produce json
// @Param schedule_id path string true "Schedule ID"
// @Success 200 {array} models.TwReminder
// @Router /api/v1/reminder/schedule/{schedule_id} [get]
func (h ReminderHandler) GetRemindersByScheduleId(ctx *fiber.Ctx) error {
	scheduleId := ctx.Params("schedule_id")
	var reminders []models.TwReminder
	if result := h.DB.
		Where("schedule_id = ?", scheduleId).
		Where("deleted_at IS NULL").
		Find(&reminders); result.Error != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}
	return ctx.JSON(reminders)
}

// updateReminder godoc
// @Summary Update a reminder
// @Description Update a reminder
// @Tags reminder
// @Accept json
// @Produce json
// @Param reminder_id path string true "Reminder ID"
// @Param reminder body models.TwReminder true "Reminder"
// @Success 200 {object} models.TwReminder
// @Router /api/v1/reminder/{reminder_id} [put]
func (h ReminderHandler) UpdateReminder(ctx *fiber.Ctx) error {
	id := ctx.Params("reminder_id")
	var reminder models.TwReminder
	if result := h.DB.
		Where("id = ?", id).
		Where("deleted_at IS NULL").
		First(&reminder); result.Error != nil {
		return ctx.Status(fiber.StatusNotFound).SendString(result.Error.Error())
	}
	if err := ctx.BodyParser(&reminder); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	if result := h.DB.Model(reminder).Update("updated_at", gorm.Expr("NOW()")); result.Error != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}

	if result := h.DB.Omit("deleted_at", "created_at").Save(reminder); result.Error != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}
	return ctx.JSON(reminder)
}

func (h ReminderHandler) DeleteReminder(ctx *fiber.Ctx) error {

	id := ctx.Params("reminder_id")
	var reminder models.TwReminder
	if result := h.DB.
		Where("id = ?", id).
		Where("deleted_at IS NULL").
		First(&reminder); result.Error != nil {
		return ctx.Status(fiber.StatusNotFound).SendString(result.Error.Error())
	}
	if result := h.DB.Model(reminder).Update("deleted_at", gorm.Expr("NOW()")); result.Error != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}

	if result := h.DB.Omit("created_at", "updated_at").Save(reminder); result.Error != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Reminder deleted successfully",
	})
}
