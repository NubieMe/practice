package internal

import (
	"practice/config"
	todoRouter "practice/internal/todo/router"
	userRouter "practice/internal/user/router"
	"practice/pkg/logger"

	"github.com/gofiber/fiber/v2"
)

func MainRoutes(f *fiber.App, db *config.DB, logger *logger.Logger) {
	api := f.Group("/api")

	userRouter.Route(api, db, logger)
	todoRouter.Route(api, db, logger)
}
