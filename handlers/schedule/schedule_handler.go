package schedule

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/timewise-team/timewise-models/dtos/core_dtos"
	"github.com/timewise-team/timewise-models/models"
	"gorm.io/gorm"
	"strconv"
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
			StartTime:         *schedule.StartTime,
			EndTime:           *schedule.EndTime,
			Location:          schedule.Location,
			CreatedBy:         schedule.CreatedBy,
			CreatedAt:         *schedule.CreatedAt,
			UpdatedAt:         *schedule.UpdatedAt,
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
			StartTime:         *schedule.StartTime,
			EndTime:           *schedule.EndTime,
			Location:          schedule.Location,
			CreatedBy:         schedule.CreatedBy,
			CreatedAt:         *schedule.CreatedAt,
			UpdatedAt:         *schedule.UpdatedAt,
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

	var startTime, endTime, createdAt, updatedAt time.Time

	if schedule.StartTime != nil {
		startTime = *schedule.StartTime
	}
	if schedule.EndTime != nil {
		endTime = *schedule.EndTime
	}
	if schedule.CreatedAt != nil {
		createdAt = *schedule.CreatedAt
	}
	if schedule.UpdatedAt != nil {
		updatedAt = *schedule.UpdatedAt
	}

	scheduleDTO := core_dtos.TwScheduleResponse{
		ID:                int(schedule.ID),
		WorkspaceID:       schedule.WorkspaceId,
		BoardColumnID:     schedule.BoardColumnId,
		Title:             schedule.Title,
		Description:       schedule.Description,
		StartTime:         startTime,
		EndTime:           endTime,
		Location:          schedule.Location,
		CreatedBy:         schedule.CreatedBy,
		CreatedAt:         createdAt,
		UpdatedAt:         updatedAt,
		Status:            schedule.Status,
		AllDay:            schedule.AllDay,
		Visibility:        schedule.Visibility,
		VideoTranscript:   &schedule.VideoTranscript,
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

	var existingSchedules []models.TwSchedule
	if err := h.DB.Where("board_column_id = ?", *scheduleDTO.BoardColumnID).Find(&existingSchedules).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	now := time.Now()

	schedule := models.TwSchedule{
		WorkspaceId:   *scheduleDTO.WorkspaceID,
		BoardColumnId: *scheduleDTO.BoardColumnID,
		Title:         *scheduleDTO.Title,
		//Description:   *scheduleDTO.Description,
		//StartTime:         *scheduleDTO.StartTime,
		//EndTime:           *scheduleDTO.EndTime,
		//Location:          *scheduleDTO.Location,
		CreatedBy: *scheduleDTO.WorkspaceUserID,
		CreatedAt: &now,
		UpdatedAt: &now,
		//Status:            *scheduleDTO.Status,
		//AllDay:            *scheduleDTO.AllDay,
		//Visibility:        *scheduleDTO.Visibility,
		//ExtraData:         *scheduleDTO.ExtraData,
		//IsDeleted:         false,
		//RecurrencePattern: *scheduleDTO.RecurrencePattern,
		Position: len(existingSchedules) + 1,
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

	now = time.Now()
	newScheduleParticipant := models.TwScheduleParticipant{
		ScheduleId:       schedule.ID,
		WorkspaceUserId:  *scheduleDTO.WorkspaceUserID,
		AssignAt:         &now,
		AssignBy:         *scheduleDTO.WorkspaceUserID,
		Status:           "participant",
		ResponseTime:     &now,
		InvitationSentAt: &now,
		InvitationStatus: "joined",
	}

	if result := h.DB.Create(&newScheduleParticipant); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(core_dtos.TwCreateShecduleResponse{
		ID:            schedule.ID,
		WorkspaceID:   schedule.WorkspaceId,
		BoardColumnID: schedule.BoardColumnId,
		Title:         schedule.Title,
		Position:      schedule.Position,
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
	workspaceUserIdStr := c.Params("workspace_user_id")
	workspaceUserId, err := strconv.Atoi(workspaceUserIdStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid workspace_user_id")
	}

	if err := h.DB.Where("id = ?", scheduleId).First(&schedule).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).SendString("Schedule not found")
		}
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	// Tạo danh sách các log khi trường được cập nhật
	var logs []models.TwScheduleLog

	// Hàm phụ: Kiểm tra và ghi log nếu có thay đổi
	checkAndLog := func(field, oldValue, newValue string) {
		if oldValue != newValue {
			logs = append(logs, models.TwScheduleLog{
				ScheduleId:      schedule.ID,
				WorkspaceUserId: workspaceUserId,
				Action:          "update schedule",
				FieldChanged:    field,
				OldValue:        oldValue,
				NewValue:        newValue,
			})
		}
	}

	// Kiểm tra và cập nhật các trường nếu có thay đổi
	if scheduleDTO.BoardColumnID != nil {
		checkAndLog("board_column_id", strconv.Itoa(schedule.BoardColumnId), strconv.Itoa(*scheduleDTO.BoardColumnID))
		schedule.BoardColumnId = *scheduleDTO.BoardColumnID
	}
	if scheduleDTO.Title != nil {
		checkAndLog("title", schedule.Title, *scheduleDTO.Title)
		schedule.Title = *scheduleDTO.Title
	}
	if scheduleDTO.Description != nil {
		checkAndLog("description", schedule.Description, *scheduleDTO.Description)
		schedule.Description = *scheduleDTO.Description
	}
	if scheduleDTO.StartTime != nil {
		checkAndLog("start_time", schedule.StartTime.String(), scheduleDTO.StartTime.String())
		schedule.StartTime = scheduleDTO.StartTime
	}
	if scheduleDTO.EndTime != nil {
		checkAndLog("end_time", schedule.EndTime.String(), scheduleDTO.EndTime.String())
		schedule.EndTime = scheduleDTO.EndTime
	}
	if scheduleDTO.Location != nil {
		checkAndLog("location", schedule.Location, *scheduleDTO.Location)
		schedule.Location = *scheduleDTO.Location
	}
	if scheduleDTO.Status != nil {
		checkAndLog("status", schedule.Status, *scheduleDTO.Status)
		schedule.Status = *scheduleDTO.Status
	}
	if scheduleDTO.AllDay != nil {
		checkAndLog("all_day", strconv.FormatBool(schedule.AllDay), strconv.FormatBool(*scheduleDTO.AllDay))
		schedule.AllDay = *scheduleDTO.AllDay
	}
	if scheduleDTO.Visibility != nil {
		checkAndLog("visibility", schedule.Visibility, *scheduleDTO.Visibility)
		schedule.Visibility = *scheduleDTO.Visibility
	}
	if scheduleDTO.ExtraData != nil {
		checkAndLog("extra_data", schedule.ExtraData, *scheduleDTO.ExtraData)
		schedule.ExtraData = *scheduleDTO.ExtraData
	}
	if scheduleDTO.IsDeleted != nil {
		checkAndLog("is_deleted", strconv.FormatBool(schedule.IsDeleted), strconv.FormatBool(*scheduleDTO.IsDeleted))
		schedule.IsDeleted = *scheduleDTO.IsDeleted
	}
	if scheduleDTO.RecurrencePattern != nil {
		checkAndLog("recurrence_pattern", schedule.RecurrencePattern, *scheduleDTO.RecurrencePattern)
		schedule.RecurrencePattern = *scheduleDTO.RecurrencePattern
	}

	// **Kiểm tra và cập nhật trường position và priority nếu có thay đổi**
	if scheduleDTO.Position != nil {
		checkAndLog("position", strconv.Itoa(schedule.Position), strconv.Itoa(*scheduleDTO.Position))
		schedule.Position = *scheduleDTO.Position
	}
	if scheduleDTO.Priority != nil {
		checkAndLog("priority", schedule.Priority, *scheduleDTO.Priority)
		schedule.Priority = *scheduleDTO.Priority
	}

	// Update timestamp
	now := time.Now()
	schedule.UpdatedAt = &now

	// Lưu schedule đã cập nhật
	if result := h.DB.Omit("deleted_at").Save(&schedule); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}

	// Thêm các log vào cơ sở dữ liệu
	if len(logs) > 0 {
		if result := h.DB.Create(&logs); result.Error != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
		}
	}

	// Trả về kết quả cập nhật thành công
	return c.JSON(core_dtos.TwUpdateScheduleResponse{
		ID:                schedule.ID,
		WorkspaceID:       schedule.WorkspaceId,
		BoardColumnID:     schedule.BoardColumnId,
		Title:             schedule.Title,
		Description:       schedule.Description,
		StartTime:         *schedule.StartTime,
		EndTime:           *schedule.EndTime,
		Location:          schedule.Location,
		CreatedBy:         schedule.CreatedBy,
		CreatedAt:         *schedule.CreatedAt,
		UpdatedAt:         *schedule.UpdatedAt,
		Status:            schedule.Status,
		AllDay:            schedule.AllDay,
		Visibility:        schedule.Visibility,
		ExtraData:         schedule.ExtraData,
		IsDeleted:         schedule.IsDeleted,
		RecurrencePattern: schedule.RecurrencePattern,
		Position:          schedule.Position,
		Priority:          schedule.Priority,
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
	workspaceUserIdStr := c.Params("workspace_user_id")
	workspaceUserId, err := strconv.Atoi(workspaceUserIdStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid workspace_user_id")
	}

	var schedule models.TwSchedule
	if err := h.DB.Where("id = ?", scheduleId).First(&schedule).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).SendString("Schedule not found")
		}
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	now := time.Now()

	schedule.IsDeleted = true
	schedule.UpdatedAt = &now
	schedule.DeletedAt = &now

	if result := h.DB.Omit("start_time,end_time").Save(&schedule); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}

	var schedulesToUpdate []models.TwSchedule
	if err := h.DB.Where("board_column_id = ? AND position > ? AND is_deleted != 1", schedule.BoardColumnId, schedule.Position).
		Order("position ASC").Find(&schedulesToUpdate).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	for i := range schedulesToUpdate {
		schedulesToUpdate[i].Position -= 1
		if err := h.DB.Omit("deleted_at,start_time,end_time").Save(&schedulesToUpdate[i]).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}
	}

	newScheduleLog := models.TwScheduleLog{
		ScheduleId:      schedule.ID,
		WorkspaceUserId: workspaceUserId,
		Action:          "delete schedule",
	}

	if result := h.DB.Create(&newScheduleLog); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}

	return c.SendStatus(fiber.StatusOK)
}

func (h *ScheduleHandler) GetSchedulesByBoardColumn(c *fiber.Ctx) error {
	boardColumnID := c.Params("board_column_id")
	if boardColumnID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid board column ID",
		})
	}
	var schedules []models.TwSchedule
	if result := h.DB.Where("board_column_id = ?", boardColumnID).Find(&schedules); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": result.Error.Error(),
		})
	}
	if schedules == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to get schedules",
		})
	}
	return c.JSON(schedules)
}

func (h *ScheduleHandler) GetSchedulesByWorkspace(c *fiber.Ctx) error {
	workspaceID := c.Params("workspace_id")
	if workspaceID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid workspace ID",
		})
	}
	var schedules []models.TwSchedule
	if result := h.DB.Where("workspace_id = ?", workspaceID).Find(&schedules); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": result.Error.Error(),
		})
	}
	if schedules == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to get schedules",
		})
	}
	return c.JSON(schedules)
}

func (h *ScheduleHandler) getSchedulesByBoardColumn(c *fiber.Ctx) error {
	boardColumnID := c.Params("board_column_id")
	workspaceID := c.Params("workspace_id")
	if boardColumnID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid board column ID",
		})
	}
	var schedules []models.TwSchedule
	if result := h.DB.Where("board_column_id = ? and workspace_id = ?", boardColumnID, workspaceID).Find(&schedules); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": result.Error.Error(),
		})
	}
	if schedules == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to get schedules",
		})
	}
	return c.JSON(schedules)
}

func (h *ScheduleHandler) UpdateTranscriptBySchedule(ctx *fiber.Ctx) error {
	// get an api_key from params
	apiKey := ctx.Get("x_api_key")
	if apiKey != "667qwsrUlyVa" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	// Parse schedule_id from params
	scheduleId := ctx.Params("schedule_id")
	if scheduleId == "" {
		return ctx.Status(fiber.StatusBadRequest).SendString("Schedule ID is required")
	}

	// Parse request body
	var scheduleDTO core_dtos.TwUpdateScheduleRequest
	if err := ctx.BodyParser(&scheduleDTO); err != nil {
		return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	// Fetch the schedule from the database
	var schedule models.TwSchedule
	if err := h.DB.Where("id = ?", scheduleId).First(&schedule).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ctx.Status(fiber.StatusNotFound).SendString("Schedule not found")
		}
		return ctx.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	// Update the VideoTranscript field if provided
	if scheduleDTO.VideoTranscript != nil {
		schedule.VideoTranscript = *scheduleDTO.VideoTranscript

		now := time.Now()
		schedule.UpdatedAt = &now
	}

	// Save the updated schedule back to the database
	if result := h.DB.Save(&schedule); result.Error != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}

	// Return the updated schedule in the response
	return ctx.JSON("Updated successfully")
}
