package feature

import (
	"dbms/handlers/user"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func RegisterHandlerV1(db *gorm.DB) *fiber.App {
	router := fiber.New()
	v1 := router.Group("/dbms/v1")

	user.RegisterUserHandler(v1.Group("/user"), db)

	return router
}
