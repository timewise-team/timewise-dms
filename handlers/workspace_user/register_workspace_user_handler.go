package workspace_user

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func RegisterWorkspaceUserHandler(router fiber.Router, db *gorm.DB) {
	workspaceUserHandler := WorkspaceUserHandler{
		Router: router,
		DB:     db,
	}

	// Register all endpoints here
	router.Get("/", workspaceUserHandler.getWorkspaceUsers)
	router.Get("/:workspace_user_id", workspaceUserHandler.getWorkspaceUserById)
	router.Delete("/:workspace_user_id", workspaceUserHandler.removeWorkspaceUserById)
	router.Post("/", workspaceUserHandler.createWorkspaceUser)
	router.Put("/", workspaceUserHandler.updateWorkspaceUser)
	router.Get("/workspace/:workspace_id", workspaceUserHandler.getWorkspaceUsersByWorkspaceId)
	router.Get("/user/:user_id", workspaceUserHandler.getWorkspaceUsersByUserId)
	router.Get("/workspace_key/:workspace_key", workspaceUserHandler.getWorkspaceUsersByWorkspaceKey)
	router.Get("/status/:status", workspaceUserHandler.getWorkspaceUsersByStatus)
	router.Get("/is_active/:is_active", workspaceUserHandler.getWorkspaceUsersByIsActive)
	router.Get("/email/:email/workspace/:workspace_id", workspaceUserHandler.getWorkspaceUserByEmailAndWorkspace)
}
