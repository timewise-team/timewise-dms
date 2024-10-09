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

// FilterSchedules godoc
// @Summary Filter schedule
// @Description Filter schedules
// @Tags schedule
// @Accept json
// @Produce json
// @Param workspace_id query int false "Workspace ID"
// @Param board_column_id query int false "Board Column ID"
// @Param title query string false "Title of the schedule (searches with LIKE)"
// @Param start_time query string false "Start time of the schedule (ISO8601 format, filter by schedules starting after this date)"
// @Param end_time query string false "End time of the schedule (ISO8601 format, filter by schedules ending before this date)"
// @Param location query string false "Location of the schedule (searches with LIKE)"
// @Param created_by query int false "User ID of the creator"
// @Param status query string false "Status of the schedule"
// @Param is_deleted query bool false "Filter by deleted schedules"
// @Param assigned_to query int false "User ID assigned to the schedule"
// @Success 200 {array} core_dtos.TwScheduleResponse "Filtered list of schedules"
// @Failure 400 {object} fiber.Error "Invalid query parameters"
// @Failure 500 {object} fiber.Error "Internal Server Error"
// @Router /dbms/v1/schedule/schedules/filter [get]
func (h *ScheduleHandler) FilterSchedules(c *fiber.Ctx) error {
	var schedules []models.TwSchedule

	query := h.DB.Model(&models.TwSchedule{})

	workspaceID := c.Query("workspace_id")
	boardColumnID := c.Query("board_column_id")
	title := c.Query("title")
	startTime := c.Query("start_time")
	endTime := c.Query("end_time")
	location := c.Query("location")
	createdBy := c.Query("created_by")
	status := c.Query("status")
	isDeleted := c.Query("is_deleted")
	assignedTo := c.Query("assigned_to")

	if workspaceID != "" {
		query = query.Where("workspace_id = ?", workspaceID)
	}

	if boardColumnID != "" {
		query = query.Where("board_column_id = ?", boardColumnID)
	}

	if title != "" {
		query = query.Where("title LIKE ?", "%"+title+"%")
	}

	if startTime != "" {
		parsedStartTime, err := time.Parse(time.RFC3339, startTime)
		if err == nil {
			query = query.Where("start_time >= ?", parsedStartTime)
		}
	}

	if endTime != "" {
		parsedEndTime, err := time.Parse(time.RFC3339, endTime)
		if err == nil {
			query = query.Where("end_time <= ?", parsedEndTime)
		}
	}

	if location != "" {
		query = query.Where("location LIKE ?", "%"+location+"%")
	}

	if createdBy != "" {
		query = query.Where("created_by = ?", createdBy)
	}

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if isDeleted != "" {
		if isDeleted == "true" {
			query = query.Where("is_deleted = ?", 1)
		} else if isDeleted == "false" {
			query = query.Where("is_deleted = ?", 0)
		} else {
			return c.Status(fiber.StatusBadRequest).SendString("Invalid value for is_deleted. Must be 'true' or 'false'")
		}
	}

	if assignedTo != "" {
		query = query.Where("assigned_to @> ?", "{"+assignedTo+"}")
	}

	if result := query.Find(&schedules); result.Error != nil {
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
		})
	}

	return c.JSON(scheduleDTOs)
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
		CreatedBy:         *scheduleDTO.WorkspaceUserID,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
		Status:            *scheduleDTO.Status,
		AllDay:            *scheduleDTO.AllDay,
		Visibility:        *scheduleDTO.Visibility,
		ExtraData:         *scheduleDTO.ExtraData,
		IsDeleted:         false,
		RecurrencePattern: *scheduleDTO.RecurrencePattern,
	}

	if result := h.DB.Create(&schedule); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}

	newScheduleLog := models.TwScheduleLog{
		ScheduleId:      schedule.ID,
		WorkspaceUserId: *scheduleDTO.WorkspaceUserID,
		Action:          "create schedule",
	}

	if result := h.DB.Create(&newScheduleLog); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}

	newScheduleParticipant := models.TwScheduleParticipant{
		ScheduleId:       schedule.ID,
		WorkspaceUserId:  *scheduleDTO.WorkspaceUserID,
		AssignAt:         time.Now(),
		AssignBy:         *scheduleDTO.WorkspaceUserID,
		Status:           "participant",
		ResponseTime:     time.Now(),
		InvitationSentAt: time.Now(),
		InvitationStatus: "joined",
	}

	if result := h.DB.Create(&newScheduleParticipant); result.Error != nil {
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

	// Update the timestamp
	schedule.UpdatedAt = time.Now()

	if result := h.DB.Omit("deleted_at").Save(&schedule); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}

	newScheduleLog := models.TwScheduleLog{
		ScheduleId:      schedule.ID,
		WorkspaceUserId: schedule.WorkspaceId,
		Action:          "update schedule",
	}

	if result := h.DB.Create(&newScheduleLog); result.Error != nil {
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
	schedule.DeletedAt = time.Now()

	if result := h.DB.Save(&schedule); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}

	return c.SendStatus(fiber.StatusNoContent)
}
