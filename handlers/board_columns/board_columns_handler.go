package board_columns

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/timewise-team/timewise-models/dtos/core_dtos/board_columns_dtos"
	"github.com/timewise-team/timewise-models/models"
	"gorm.io/gorm"
)

// getBoardColumnsByWorkspace godoc
// @Summary Get board columns by workspace
// @Description Get board columns by workspace
// @Tags board_columns
// @Accept json
// @Produce json
// @Param workspace_id path int true "Workspace ID"
// @Success 200 {array} models.TwBoardColumn
// @Router /dbms/v1/workspace/{workspace_id}/board_columns [get]
func (h *BoardColumnsHandler) getBoardColumnsByWorkspace(c *fiber.Ctx) error {
	// Parse the request
	workspaceID := c.Params("workspace_id")
	if workspaceID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid workspace ID",
		})
	}
	var boardColumns []models.TwBoardColumn
	// Get the board columns
	if result := h.DB.Where("workspace_id = ?", workspaceID).Find(&boardColumns); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": result.Error.Error(),
		})
	}
	if boardColumns == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to get board columns",
		})
	}
	// Return the response
	return c.JSON(boardColumns)
}

// getBoardColumnById godoc
// @Summary Get board column by ID
// @Description Get board column by ID
// @Tags board_columns
// @Accept json
// @Produce json
// @Param id path int true "Board column ID"
// @Success 200 {object} models.TwBoardColumn
// @Router /dbms/v1/board_columns/{id} [get]
func (h *BoardColumnsHandler) getBoardColumnById(c *fiber.Ctx) error {
	var boardColumn models.TwBoardColumn
	boardColumnId := c.Params("id")

	if err := h.DB.Where("id = ?", boardColumnId).First(&boardColumn).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).SendString("BoardColumn not found")
		}
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	return c.JSON(boardColumn)
}

// deleteBoardColumn godoc
// @Summary Delete board column
// @Description Delete board column
// @Tags board_columns
// @Accept json
// @Produce json
// @Param id path int true "Board column ID"
// @Success 204
// @Router /dbms/v1/board_columns/{id} [delete]
func (h *BoardColumnsHandler) deleteBoardColumn(c *fiber.Ctx) error {
	boardColumnId := c.Params("id")
	var boardColumn models.TwBoardColumn
	if result := h.DB.Where("id = ?", boardColumnId).Delete(&boardColumn); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}
	return c.SendStatus(fiber.StatusNoContent)
}

// updateBoardColumn godoc
// @Summary Update board column
// @Description Update board column
// @Tags board_columns
// @Accept json
// @Produce json
// @Param id path int true "Board column ID"
// @Param body body models.TwBoardColumn true "Update board column request"
// @Success 200 {object} models.TwBoardColumn
// @Router /dbms/v1/board_columns/{id} [put]
func (h *BoardColumnsHandler) updateBoardColumn(c *fiber.Ctx) error {
	boardColumnId := c.Params("id")
	var boardColumn models.TwBoardColumn
	if err := h.DB.Where("id = ?", boardColumnId).First(&boardColumn).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).SendString("BoardColumn not found")
		}
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	// Parse the request
	var updateBoardColumnRequest models.TwBoardColumn
	if err := c.BodyParser(&updateBoardColumnRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	// Update the board column
	if result := h.DB.Model(&boardColumn).Updates(updateBoardColumnRequest); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}
	return c.JSON(boardColumn)
}

// getBoardColumnField godoc
// @Summary Get board column field
// @Description Get board column field
// @Tags board_columns
// @Accept json
// @Produce json
// @Param id path int true "Board column ID"
// @Param field path string true "Field"
// @Success 200 {object} models.TwBoardColumn
// @Router /dbms/v1/board_columns/{id}/{field} [get]
func (h *BoardColumnsHandler) getBoardColumnField(c *fiber.Ctx) error {
	field := c.Params("field")
	boardColumnId := c.Params("id")
	if boardColumnId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid board column ID",
		})
	}
	if field == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid field",
		})
	}
	if field != "name" && field != "position" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid field",
		})
	}
	if field == "name" {
		var boardColumn models.TwBoardColumn
		if err := h.DB.Where("id = ?", boardColumnId).First(&boardColumn).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return c.Status(fiber.StatusNotFound).SendString("BoardColumn not found")
			}
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}
		return c.JSON(fiber.Map{
			"name": boardColumn.Name,
		})
	}
	if field == "position" {
		var boardColumn models.TwBoardColumn
		if err := h.DB.Where("id = ?", boardColumnId).First(&boardColumn).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return c.Status(fiber.StatusNotFound).SendString("BoardColumn not found")
			}
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}
		return c.JSON(fiber.Map{
			"position": boardColumn.Position,
		})
	}
	return c.Status(fiber.StatusInternalServerError).SendString("Invalid field")
}

// updateBoardColumnField godoc
// @Summary Update board column field
// @Description Update board column field
// @Tags board_columns
// @Accept json
// @Produce json
// @Param id path int true "Board column ID"
// @Param field path string true "Field"
// @Param body body models.TwBoardColumn true "Update board column request"
// @Success 200 {object} models.TwBoardColumn
// @Router /dbms/v1/board_columns/{id}/{field} [put]
func (h *BoardColumnsHandler) updateBoardColumnField(c *fiber.Ctx) error {
	field := c.Params("field")
	boardColumnId := c.Params("id")
	if boardColumnId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid board column ID",
		})
	}
	if field == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid field",
		})
	}
	if field != "name" && field != "position" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid field",
		})
	}
	if field == "name" {
		var boardColumn models.TwBoardColumn
		if err := h.DB.Where("id = ?", boardColumnId).First(&boardColumn).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return c.Status(fiber.StatusNotFound).SendString("BoardColumn not found")
			}
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}
		var updateBoardColumnRequest models.TwBoardColumn
		if err := c.BodyParser(&updateBoardColumnRequest); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": err.Error(),
			})
		}
		if result := h.DB.Model(&boardColumn).Update("name", updateBoardColumnRequest.Name); result.Error != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
		}
		return c.JSON(updateBoardColumnRequest)

	}
	if field == "position" {
		var boardColumn models.TwBoardColumn
		if err := h.DB.Where("id = ?", boardColumnId).First(&boardColumn).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return c.Status(fiber.StatusNotFound).SendString("BoardColumn not found")
			}
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}
		var updateBoardColumnRequest models.TwBoardColumn
		if err := c.BodyParser(&updateBoardColumnRequest); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": err.Error(),
			})
		}
		if result := h.DB.Model(&boardColumn).Update("position", updateBoardColumnRequest.Position); result.Error != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
		}
		return c.JSON(fiber.Map{
			"position": boardColumn.Position,
		})
	}
	return c.Status(fiber.StatusInternalServerError).SendString("Invalid field")
}

// createBoardColumn godoc
// @Summary Create board column
// @Description Create board column
// @Tags board_columns
// @Accept json
// @Produce json
// @Param body body board_columns_dtos.BoardColumnsRequest true "Create board column request"
// @Success 200 {object} models.TwBoardColumn
// @Router /dbms/v1/board_columns [post]
func (h *BoardColumnsHandler) createBoardColumn(c *fiber.Ctx) error {
	// Parse the request
	var createBoardColumnRequest board_columns_dtos.BoardColumnsRequest
	if err := c.BodyParser(&createBoardColumnRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	var boardColumn = models.TwBoardColumn{
		Name:        createBoardColumnRequest.Name,
		Position:    createBoardColumnRequest.Position,
		WorkspaceId: createBoardColumnRequest.WorkspaceId,
	}
	// Create the board column
	if result := h.DB.Create(&boardColumn); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}
	return c.JSON(boardColumn)
}
