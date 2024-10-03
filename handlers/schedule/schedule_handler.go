package schedule

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/timewise-team/timewise-models/dtos/core_dtos"
	"github.com/timewise-team/timewise-models/models"
	"gorm.io/gorm"
	"time"
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
// @Success 200 {array} core_dtos.TwScheduleResponse
// @Router /dbms/v1/schedule [get]
func (h *ScheduleHandler) GetSchedules(c *fiber.Ctx) error {
	var schedules []models.TwSchedule
	if result := h.DB.Find(&schedules); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}

	var scheduleDTOs []core_dtos.TwScheduleResponse
	for _, schedule := range schedules {
		scheduleDTOs = append(scheduleDTOs, core_dtos.TwScheduleResponse{
			ID:                int(schedule.ID),
			WorkspaceID:       schedule.WorkspaceId,
			BoardColumnID:     schedule.BoardColumnId,
			Title:             schedule.Title,
			Description:       schedule.Description,
			StartTime:         schedule.StartTime,
			EndTime:           schedule.EndTime,
			Location:          schedule.Location,
			CreatedBy:         schedule.CreatedBy,
			CreatedAt:         schedule.CreatedAt,
			UpdatedAt:         schedule.UpdatedAt,
			Status:            schedule.Status,
			AllDay:            schedule.AllDay,
			Visibility:        schedule.Visibility,
			ExtraData:         schedule.ExtraData,
			IsDeleted:         schedule.IsDeleted,
			RecurrencePattern: schedule.RecurrencePattern,
			//AssignedTo:        []int{schedule.AssignedTo},
		})
	}

	return c.JSON(scheduleDTOs)
}

// GetScheduleById godoc
// @Summary Get schedule by ID
// @Description Get schedule by ID
// @Tags schedule
// @Accept json
// @Produce json
// @Param schedule_id path int true "Schedule ID"
// @Success 200 {object} core_dtos.TwScheduleResponse
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

	scheduleDTO := core_dtos.TwScheduleResponse{
		ID:                int(schedule.ID),
		WorkspaceID:       schedule.WorkspaceId,
		BoardColumnID:     schedule.BoardColumnId,
		Title:             schedule.Title,
		Description:       schedule.Description,
		StartTime:         schedule.StartTime,
		EndTime:           schedule.EndTime,
		Location:          schedule.Location,
		CreatedBy:         schedule.CreatedBy,
		CreatedAt:         schedule.CreatedAt,
		UpdatedAt:         schedule.UpdatedAt,
		Status:            schedule.Status,
		AllDay:            schedule.AllDay,
		Visibility:        schedule.Visibility,
		ExtraData:         schedule.ExtraData,
		IsDeleted:         schedule.IsDeleted,
		RecurrencePattern: schedule.RecurrencePattern,
		//AssignedTo:        []int{schedule.AssignedTo},
	}

	return c.JSON(scheduleDTO)
}

// CreateSchedule godoc
// @Summary Create a new schedule
// @Description Create a new schedule
// @Tags schedule
// @Accept json
// @Produce json
// @Param schedule body core_dtos.TwCreateScheduleRequest true "Schedule"
// @Success 201 {object} core_dtos.TwCreateShecduleResponse
// @Router /dbms/v1/schedule [post]
func (h *ScheduleHandler) CreateSchedule(c *fiber.Ctx) error {
	var scheduleDTO core_dtos.TwCreateScheduleRequest
	if err := c.BodyParser(&scheduleDTO); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	schedule := models.TwSchedule{
		WorkspaceId:       *scheduleDTO.WorkspaceID,
		BoardColumnId:     *scheduleDTO.BoardColumnID,
		Title:             *scheduleDTO.Title,
		Description:       *scheduleDTO.Description,
		StartTime:         *scheduleDTO.StartTime,
		EndTime:           *scheduleDTO.EndTime,
		Location:          *scheduleDTO.Location,
		CreatedBy:         scheduleDTO.CreatedBy,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
		Status:            *scheduleDTO.Status,
		AllDay:            *scheduleDTO.AllDay,
		Visibility:        *scheduleDTO.Visibility,
		ExtraData:         *scheduleDTO.ExtraData,
		IsDeleted:         false,
		RecurrencePattern: *scheduleDTO.RecurrencePattern,
		//AssignedTo:        *scheduleDTO.AssignedTo,
	}

	if result := h.DB.Create(&schedule); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(core_dtos.TwCreateShecduleResponse{
		ID:                schedule.ID,
		WorkspaceID:       schedule.WorkspaceId,
		BoardColumnID:     schedule.BoardColumnId,
		Title:             schedule.Title,
		Description:       schedule.Description,
		StartTime:         schedule.StartTime,
		EndTime:           schedule.EndTime,
		Location:          schedule.Location,
		CreatedBy:         schedule.CreatedBy,
		CreatedAt:         schedule.CreatedAt,
		UpdatedAt:         schedule.UpdatedAt,
		Status:            schedule.Status,
		AllDay:            schedule.AllDay,
		Visibility:        schedule.Visibility,
		ExtraData:         schedule.ExtraData,
		IsDeleted:         schedule.IsDeleted,
		RecurrencePattern: schedule.RecurrencePattern,
		//AssignedTo:        schedule.AssignedTo,
	})
}

// UpdateSchedule godoc
// @Summary Update an existing schedule
// @Description Update an existing schedule
// @Tags schedule
// @Accept json
// @Produce json
// @Param schedule_id path int true "Schedule ID"
// @Param schedule body core_dtos.TwUpdateScheduleRequest true "Schedule"
// @Success 200 {object} core_dtos.TwUpdateScheduleResponse
// @Router /dbms/v1/schedule/{schedule_id} [put]
func (h *ScheduleHandler) UpdateSchedule(c *fiber.Ctx) error {
	var scheduleDTO core_dtos.TwUpdateScheduleRequest
	if err := c.BodyParser(&scheduleDTO); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	var schedule models.TwSchedule
	scheduleId := c.Params("schedule_id")

	if err := h.DB.Where("id = ?", scheduleId).First(&schedule).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).SendString("Schedule not found")
		}
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	// Update the fields if they are provided (not nil)
	if scheduleDTO.WorkspaceID != nil {
		schedule.WorkspaceId = *scheduleDTO.WorkspaceID
	}
	if scheduleDTO.BoardColumnID != nil {
		schedule.BoardColumnId = *scheduleDTO.BoardColumnID
	}
	if scheduleDTO.Title != nil {
		schedule.Title = *scheduleDTO.Title
	}
	if scheduleDTO.Description != nil {
		schedule.Description = *scheduleDTO.Description
	}
	if scheduleDTO.StartTime != nil {
		schedule.StartTime = *scheduleDTO.StartTime
	}
	if scheduleDTO.EndTime != nil {
		schedule.EndTime = *scheduleDTO.EndTime
	}
	if scheduleDTO.Location != nil {
		schedule.Location = *scheduleDTO.Location
	}
	if scheduleDTO.Status != nil {
		schedule.Status = *scheduleDTO.Status
	}
	if scheduleDTO.AllDay != nil {
		schedule.AllDay = *scheduleDTO.AllDay
	}
	if scheduleDTO.Visibility != nil {
		schedule.Visibility = *scheduleDTO.Visibility
	}
	if scheduleDTO.ExtraData != nil {
		schedule.ExtraData = *scheduleDTO.ExtraData
	}
	if scheduleDTO.IsDeleted != nil {
		schedule.IsDeleted = *scheduleDTO.IsDeleted
	}
	if scheduleDTO.RecurrencePattern != nil {
		schedule.RecurrencePattern = *scheduleDTO.RecurrencePattern
	}
	//if scheduleDTO.AssignedTo != nil {
	//	schedule.AssignedTo = *scheduleDTO.AssignedTo
	//}

	// Update the timestamp
	schedule.UpdatedAt = time.Now()

	if result := h.DB.Save(&schedule); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}

	return c.JSON(core_dtos.TwUpdateScheduleResponse{
		ID:                schedule.ID,
		WorkspaceID:       schedule.WorkspaceId,
		BoardColumnID:     schedule.BoardColumnId,
		Title:             schedule.Title,
		Description:       schedule.Description,
		StartTime:         schedule.StartTime,
		EndTime:           schedule.EndTime,
		Location:          schedule.Location,
		CreatedBy:         schedule.CreatedBy,
		CreatedAt:         schedule.CreatedAt,
		UpdatedAt:         schedule.UpdatedAt,
		Status:            schedule.Status,
		AllDay:            schedule.AllDay,
		Visibility:        schedule.Visibility,
		ExtraData:         schedule.ExtraData,
		IsDeleted:         schedule.IsDeleted,
		RecurrencePattern: schedule.RecurrencePattern,
		//AssignedTo:        schedule.AssignedTo,
	})
}

// DeleteSchedule godoc
// @Summary Delete a schedule
// @Description Delete a schedule
// @Tags schedule
// @Accept json
// @Produce json
// @Param schedule_id path int true "Schedule ID"
// @Success 204 "No Content"
// @Router /dbms/v1/schedule/{schedule_id} [delete]
func (h *ScheduleHandler) DeleteSchedule(c *fiber.Ctx) error {
	scheduleId := c.Params("schedule_id")

	var schedule models.TwSchedule
	if err := h.DB.Where("id = ?", scheduleId).First(&schedule).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).SendString("Schedule not found")
		}
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	// Soft delete by setting IsDeleted to true
	schedule.IsDeleted = true
	schedule.UpdatedAt = time.Now()

	if result := h.DB.Save(&schedule); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}

	return c.SendStatus(fiber.StatusNoContent)
}
