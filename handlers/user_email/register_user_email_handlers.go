package user_email

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type UserEmailHandler struct {
	Router fiber.Router
	DB     *gorm.DB
}

func RegisterUserEmailHandler(router fiber.Router, db *gorm.DB) {
	userEmailHandler := UserEmailHandler{
		Router: router,
		DB:     db,
	}

	router.Get("/", userEmailHandler.getUserEmails)
	router.Get("/user/:user_id", userEmailHandler.getUserEmailByUserId)
	router.Get("/email/:email", userEmailHandler.getUserEmailByEmail)
	router.Post("/", userEmailHandler.createUserEmail)
	router.Patch("/", userEmailHandler.updateUserIdInUserEmail)
	router.Delete("/:email_id", userEmailHandler.deleteUserEmail)
	router.Get("/search/:query", userEmailHandler.searchUserEmail)
}
