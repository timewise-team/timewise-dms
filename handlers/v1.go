package feature

import (
	_ "dbms/docs"
	"dbms/handlers/user"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"gorm.io/gorm"
)

// @host localhost:8080
// @BasePath /dbms/v1
func RegisterHandlerV1(db *gorm.DB) *fiber.App {
	router := fiber.New()
	v1 := router.Group("/dbms/v1")
	v1.Get("/swagger/*", swagger.HandlerDefault)
	user.RegisterUserHandler(v1.Group("/user"), db)

	return router
}
