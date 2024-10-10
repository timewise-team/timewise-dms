package board_columns

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/timewise-team/timewise-models/models"
	"gorm.io/gorm"
)

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

func (h *BoardColumnsHandler) deleteBoardColumn(c *fiber.Ctx) error {
	boardColumnId := c.Params("id")
	var boardColumn models.TwBoardColumn
	if result := h.DB.Where("id = ?", boardColumnId).Delete(&boardColumn); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}
	return c.SendStatus(fiber.StatusNoContent)
}

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

func (h *BoardColumnsHandler) createBoardColumn(c *fiber.Ctx) error {
	// Parse the request
	var createBoardColumnRequest models.TwBoardColumn
	if err := c.BodyParser(&createBoardColumnRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	// Create the board column
	if result := h.DB.Create(&createBoardColumnRequest); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}
	return c.JSON(createBoardColumnRequest)
}
