package workspace_user

//workspace_user_handler.go
import (
	"errors"
	"github.com/gofiber/fiber/v2"
	workspaceUserDtos "github.com/timewise-team/timewise-models/dtos/core_dtos/workspace_user_dtos"
	"github.com/timewise-team/timewise-models/models"

	"gorm.io/gorm"
	"net/http"
	"net/url"
)

type WorkspaceUserHandler struct {
	Router fiber.Router
	DB     *gorm.DB
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

//func (h *WorkspaceUserHandler) getWorkspaceUsersByWorkspaceId(c *fiber.Ctx) error {
//	var workspaceUsers []models.TwWorkspaceUser
//	workspaceId := c.Params("workspace_id")
//
//	if result := h.DB.Where("workspace_id = ?", workspaceId).Find(&workspaceUsers); result.Error != nil {
//		return c.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
//	}
//
//	return c.JSON(workspaceUsers)
//}

// getWorkspaceUsersByWorkspaceId godoc
// @Summary Get workspace users by workspace ID
// @Description Get workspace users by workspace ID
// @Tags workspace_user
// @Accept json
// @Produce json
// @Param workspace_id path string true "Workspace ID"
// @Success 200 {array} workspaceUserDtos.GetWorkspaceUserListResponse
// @Router /dbms/v1/workspace_user/workspace/{workspace_id} [get]
func (h *WorkspaceUserHandler) getWorkspaceUsersByWorkspaceId(c *fiber.Ctx) error {
	var workspaceUsers []workspaceUserDtos.GetWorkspaceUserListResponse
	workspaceId := c.Params("workspace_id")
	if workspaceId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "workspace_id is required",
		})
	}
	err := h.DB.Table("tw_workspace_users").
		Select("tw_workspace_users.id, tw_workspace_users.user_email_id, tw_workspace_users.workspace_id, tw_workspace_users.workspace_key,tw_workspace_users.role,  tw_workspace_users.status, tw_workspace_users.is_active, tw_workspace_users.is_verified,  tw_workspace_users.extra_data, tw_workspace_users.created_at, tw_workspace_users.updated_at, tw_workspace_users.deleted_at, tw_user_emails.email, tw_users.first_name,tw_users.last_name,tw_users.profile_picture").
		Joins("JOIN tw_user_emails ON tw_workspace_users.user_email_id= tw_user_emails.id").
		Joins("JOIN tw_users ON tw_user_emails.user_id = tw_users.id").
		Where("tw_workspace_users.workspace_id = ? and tw_users.is_verified = true and tw_users.is_active = true and tw_workspace_users.status = 'joined' and tw_workspace_users.is_active = true and tw_workspace_users.is_verified=true", workspaceId).
		Scan(&workspaceUsers).Error
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	return c.JSON(workspaceUsers)
}

func (h *WorkspaceUserHandler) getWorkspaceUsersByUserId(c *fiber.Ctx) error {
	var workspaceUsers []models.TwWorkspaceUser
	userId := c.Params("user_id")

	if result := h.DB.Where("user_id = ?", userId).Find(&workspaceUsers); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}

	return c.JSON(workspaceUsers)
}

func (h *WorkspaceUserHandler) getWorkspaceUsersByWorkspaceKey(c *fiber.Ctx) error {
	var workspaceUsers []models.TwWorkspaceUser
	workspaceKey := c.Params("workspace_key")

	if result := h.DB.Where("workspace_key = ?", workspaceKey).Find(&workspaceUsers); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}

	return c.JSON(workspaceUsers)
}

func (h *WorkspaceUserHandler) getWorkspaceUsersByStatus(c *fiber.Ctx) error {
	var workspaceUsers []models.TwWorkspaceUser
	status := c.Params("status")

	if result := h.DB.Where("status = ?", status).Find(&workspaceUsers); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}

	return c.JSON(workspaceUsers)
}

func (h *WorkspaceUserHandler) getWorkspaceUsersByIsActive(c *fiber.Ctx) error {
	var workspaceUsers []models.TwWorkspaceUser
	isActive := c.Params("is_active")

	if result := h.DB.Where("is_active = ?", isActive).Find(&workspaceUsers); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}

	return c.JSON(workspaceUsers)
}

// GET /workspaces/email/{email}
// getWorkspacesByEmail godoc
// @Summary Get workspaces by email
// @Description Get workspaces by email
// @Tags workspace
// @Accept json
// @Produce json
// @Param email path string true "Email"
// @Param workspace_id path string true "Workspace ID"
// @Success 200 {object} models.TwWorkspace
// @Router /dbms/v1/workspace_user/email/{email}/workspace/{workspace_id} [get]
func (h *WorkspaceUserHandler) getWorkspaceUserByEmailAndWorkspace(c *fiber.Ctx) error {
	workspaceId := c.Params("workspace_id")
	email := c.Params("email")
	emailFix, err1 := url.QueryUnescape(email)
	if err1 != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid email",
		})
	}
	if workspaceId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "workspace_id is required",
		})
	}
	if email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "email is required",
		})
	}
	var TWorkspaceUser models.TwWorkspaceUser

	err := h.DB.
		Table("tw_workspace_users").
		Select("tw_workspace_users.id, tw_workspace_users.created_at, tw_workspace_users.updated_at, tw_workspace_users.deleted_at, tw_workspace_users.user_email_id, tw_workspace_users.workspace_id, tw_workspace_users.workspace_key, tw_workspace_users.role, tw_workspace_users.status, tw_workspace_users.is_active,tw_workspace_users.is_verified,tw_workspace_users.extra_data").
		Joins("JOIN tw_user_emails ON tw_workspace_users.user_email_id= tw_user_emails.id").
		Where("tw_user_emails.email = ? and tw_workspace_users.workspace_id=? ", emailFix, workspaceId).
		Scan(&TWorkspaceUser).Error

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	return c.JSON(TWorkspaceUser)

}

func (h *WorkspaceUserHandler) GetWorkspaceUserInvitationList(c *fiber.Ctx) error {
	var workspaceUsers []workspaceUserDtos.GetWorkspaceUserListResponse
	workspaceId := c.Params("workspace_id")
	if workspaceId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "workspace_id is required",
		})
	}
	err := h.DB.Table("tw_workspace_users").
		Select("tw_workspace_users.id, tw_workspace_users.user_email_id, tw_workspace_users.workspace_id, tw_workspace_users.workspace_key,tw_workspace_users.role,  tw_workspace_users.status, tw_workspace_users.is_active, tw_workspace_users.is_verified,  tw_workspace_users.extra_data, tw_workspace_users.created_at, tw_workspace_users.updated_at, tw_workspace_users.deleted_at, tw_user_emails.email, tw_users.first_name,tw_users.last_name,tw_users.profile_picture").
		Joins("JOIN tw_user_emails ON tw_workspace_users.user_email_id= tw_user_emails.id").
		Joins("JOIN tw_users ON tw_user_emails.user_id = tw_users.id").
		Where("tw_workspace_users.workspace_id = ? and tw_users.is_verified = true and tw_users.is_active = true and tw_workspace_users.status != 'joined' and tw_workspace_users.is_active = true ", workspaceId).
		Scan(&workspaceUsers).Error
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	return c.JSON(workspaceUsers)
}
