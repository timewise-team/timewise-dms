package routes

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type UserHandler struct {
	Router fiber.Router
	DB     *gorm.DB
}

func RegisterUserHandler(router fiber.Router, db *gorm.DB) {
	userHandler := UserHandler{
		Router: router,
		DB:     db,
	}

	// Register all endpoints here
	router.Get("/", userHandler.getUsers)
	router.Get("/:user_id", userHandler.getUserById)
}

// GET /users
func (h *UserHandler) getUsers(c *fiber.Ctx) error {
	return c.SendString("Get all users")
}

// GET /users/{id}
func (h *UserHandler) getUserById(c *fiber.Ctx) error {
	return c.SendString("Get user by id " + c.Params("user_id"))
}
