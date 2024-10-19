package email_synced

import (
	"dbms/common"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func RegisterEmailSyncedHandler(router fiber.Router, db *gorm.DB) {
	emailSyncedHandler := &EmailSyncedHandler{
		DB: db,
	}
	common.RegisterHandler(router, db, func(handler common.Handler) {
		handler.Router.Get("/:email", emailSyncedHandler.GetAllEmailSynced)
	})
}
