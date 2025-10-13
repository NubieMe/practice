package router

import (
	"practice/config"
	"practice/internal/todo/handler"
	"practice/internal/todo/repository"
	"practice/internal/todo/usecase"
	"practice/pkg/bus"
	"practice/pkg/logger"
	"practice/pkg/middleware"
	"practice/pkg/validator"

	"github.com/gofiber/fiber/v2"
)

func Route(f fiber.Router, db *config.DB, logger *logger.Logger) {
	validator := validator.NewCustomValidator()
	event := bus.NewEventBus()

	repo := repository.NewTodoRepo(db.Instance())
	usecase := usecase.NewTodoUsecase(repo, validator, logger)
	handler := handler.NewTodoHandler(usecase, logger, event)

	todo := f.Group("/todo", middleware.JWTAuth())

	todo.Get("", handler.GetTodos)
	todo.Get("/:id", handler.GetTodo)
	todo.Post("", middleware.Upload(), handler.AddTodo)
	todo.Put("/:id", middleware.Upload(), handler.UpdateTodo)
	todo.Delete("/:id", handler.DeleteTodo)
}
