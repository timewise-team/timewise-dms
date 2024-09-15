package workspace_user

//workspace_user_handler.go
import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/timewise-team/timewise-models/models"
	"gorm.io/gorm"
)

type WorkspaceUserHandler struct {
	Router fiber.Router
	DB     *gorm.DB
}

func RegisterWorkspaceUserHandler(router fiber.Router, db *gorm.DB) {
	workspaceUserHandler := WorkspaceUserHandler{
		Router: router,
		DB:     db,
	}

	// Register all endpoints here
	router.Get("/", workspaceUserHandler.getWorkspaceUsers)
	router.Get("/:workspace_user_id", workspaceUserHandler.getWorkspaceUserById)
	router.Delete("/:workspace_user_id", workspaceUserHandler.removeWorkspaceUserById)
}

func (h *WorkspaceUserHandler) getWorkspaceUsers(c *fiber.Ctx) error {
	var workspaceUsers []models.TwWorkspaceUser
	if result := h.DB.Find(&workspaceUsers); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}

	return c.JSON(workspaceUsers)
}

func (h *WorkspaceUserHandler) getWorkspaceUserById(c *fiber.Ctx) error {
	var workspaceUser models.TwWorkspaceUser
	workspaceUserId := c.Params("workspace_user_id")

	if err := h.DB.Where("id = ?", workspaceUserId).First(&workspaceUser).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).SendString("WorkspaceUser not found")
		}
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	return c.JSON(workspaceUser)
}

func (h *WorkspaceUserHandler) removeWorkspaceUserById(c *fiber.Ctx) error {
	workspaceUserId := c.Params("workspace_user_id")
	if result := h.DB.Delete(&models.TwWorkspaceUser{}, workspaceUserId); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func (h *WorkspaceUserHandler) createWorkspaceUser(c *fiber.Ctx) error {
	workspaceUser := new(models.TwWorkspaceUser)
	if err := c.BodyParser(workspaceUser); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	if result := h.DB.Create(workspaceUser); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}

	return c.JSON(workspaceUser)
}

func (h *WorkspaceUserHandler) updateWorkspaceUser(c *fiber.Ctx) error {
	workspaceUser := new(models.TwWorkspaceUser)
	if err := c.BodyParser(workspaceUser); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	if result := h.DB.Save(workspaceUser); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}

	return c.JSON(workspaceUser)
}
func (h *WorkspaceUserHandler) getWorkspaceUsersByWorkspaceId(c *fiber.Ctx) error {
	var workspaceUsers []models.TwWorkspaceUser
	workspaceId := c.Params("workspace_id")

	if result := h.DB.Where("workspace_id = ?", workspaceId).Find(&workspaceUsers); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}

	return c.JSON(workspaceUsers)
}
