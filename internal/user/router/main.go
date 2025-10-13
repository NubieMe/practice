package router

import (
	"practice/config"
	"practice/internal/user/handler"
	"practice/internal/user/repository"
	"practice/internal/user/usecase"
	"practice/pkg/bus"
	"practice/pkg/logger"
	"practice/pkg/middleware"
	"practice/pkg/validator"

	"github.com/gofiber/fiber/v2"
)

func Route(f fiber.Router, db *config.DB, logger *logger.Logger) {
	validator := validator.NewCustomValidator()
	event := bus.NewEventBus()

	repo := repository.NewUserRepo(db.Instance())
	usecase := usecase.NewUserUsecase(repo, validator, logger)
	handler := handler.NewUserHandler(usecase, logger, event)

	f.Post("/register", handler.Register)
	f.Post("/login", handler.Login)

	user := f.Group("/user", middleware.JWTAuth())

	user.Get("/:id", handler.GetUser)
	user.Get("", handler.GetUsers)
}
