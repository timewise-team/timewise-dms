package feature

import (
	_ "dbms/docs"
	"dbms/handlers/auth"
	"dbms/handlers/recurrence_exception"
	"dbms/handlers/schedule"
	"dbms/handlers/schedule_log"
	"dbms/handlers/schedule_participant"
	"dbms/handlers/user"
	"dbms/handlers/workspace_log"
	"dbms/handlers/workspace_user"
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
	schedule_log.RegisterScheduleLogHandler(v1.Group("/schedulelog"), db)
	schedule_participant.RegisterScheduleParticipantHandler(v1.Group("/scheduleparticipant"), db)
	schedule.RegisterScheduleHandler(v1.Group("/schedule"), db)
	recurrence_exception.RegisterRecurrenceExceptionHandler(v1.Group("/recurrence_exception"), db)
	workspace_user.RegisterWorkspaceUserHandler(v1.Group("/workspace_user"), db)
	workspace_log.RegisterWorkspaceLogHandler(v1.Group("/workspace_log"), db)
	auth.RegisterAuthHandler(v1.Group("/auth"), db)
	return router
}
