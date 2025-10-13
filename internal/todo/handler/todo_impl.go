package handler

import (
	"fmt"
	"practice/internal/todo/usecase"
	"practice/models"
	"practice/pkg/bus"
	"practice/pkg/logger"
	"practice/pkg/pagination"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TodoHandlerImpl struct {
	usecase usecase.TodoUsecase
	logger  *logger.Logger
	event   *bus.EventBus
}

func NewTodoHandler(usecase usecase.TodoUsecase, logger *logger.Logger, event *bus.EventBus) TodoHandler {
	return &TodoHandlerImpl{
		usecase: usecase,
		logger:  logger,
		event:   event,
	}
}

func (h *TodoHandlerImpl) Handle(event bus.Event) {
	h.logger.Info("Todo event: %s - %v", event.Type, event.Payload)
}

func (h *TodoHandlerImpl) AddTodo(c *fiber.Ctx) error {
	request := new(models.TodoRequest)
	if err := c.BodyParser(request); err != nil {
		h.logger.Debug(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid request",
		})
	}

	userLocal := c.Locals("user")
	filenames := c.Locals("filenames")

	userMap := userLocal.(jwt.MapClaims)
	id := userMap["id"].(string)

	uid, err := uuid.Parse(id)
	if err != nil {
		h.logger.Debug(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid ID",
		})
	}

	request.UserID = uid
	request.Images = filenames.([]string)

	fmt.Println("filenames", filenames)
	fmt.Println("request", request)

	if err := h.usecase.AddTodo(c.Context(), request); err != nil {
		h.logger.Error(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "internal server error",
		})
	}

	h.event.Publish(bus.Event{
		Type:    "todo.created",
		Payload: request,
	})

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "success",
	})
}

func (h *TodoHandlerImpl) UpdateTodo(c *fiber.Ctx) error {
	id := c.Params("id", "")

	uid, err := uuid.Parse(id)
	if err != nil {
		h.logger.Debug(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid ID",
		})
	}

	// userMap := c.Locals("user").(jwt.MapClaims)
	// userID := userMap["id"].(string)

	// uuidUser, err := uuid.Parse(userID)
	// if err != nil {
	// 	h.logger.Debug(err.Error())
	// 	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
	// 		"message": "invalid ID",
	// 	})
	// }

	request := new(models.Todo)
	if err := c.BodyParser(request); err != nil {
		h.logger.Debug(err.Error())
		return c.Next()
	}
	request.ID = uid
	// request.UpdatedBy = user.ID

	if err := h.usecase.UpdateTodo(c.Context(), request); err != nil {
		h.logger.Error(err.Error())
		return c.Next()
	}

	h.event.Publish(bus.Event{
		Type:    "todo.updated",
		Payload: request,
	})

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "success",
	})
}

func (h *TodoHandlerImpl) GetTodo(c *fiber.Ctx) error {
	id := c.Params("id", "")

	uuid, err := uuid.Parse(id)
	if err != nil {
		h.logger.Debug(err.Error())
		return c.Next()
	}

	todo, err := h.usecase.GetTodo(c.Context(), uuid)
	if err != nil {
		h.logger.Error(err.Error())
		return c.Next()
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "success",
		"data":    todo,
	})
}

func (h *TodoHandlerImpl) GetTodos(c *fiber.Ctx) error {
	userMap := c.Locals("user").(jwt.MapClaims)
	uid, err := uuid.Parse(userMap["id"].(string))
	if err != nil {
		h.logger.Debug(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid ID",
		})
	}

	var pagParams models.PaginationRequest
	if err := c.QueryParser(&pagParams); err != nil {
		h.logger.Debug(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid pagination params",
		})
	}

	params := &pagination.PaginationParams{
		Page:  pagParams.Page,
		Limit: pagParams.Limit,
		Sort:  pagParams.Sort,
	}

	todos, err := h.usecase.GetTodos(c.Context(), params, uid, "todo", c.Query("todo"))
	if err != nil {
		h.logger.Error(err.Error())
		return c.Next()
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "success",
		"data":    todos,
	})
}

func (h *TodoHandlerImpl) DeleteTodo(c *fiber.Ctx) error {
	id := c.Params("id", "")

	uuid, err := uuid.Parse(id)
	if err != nil {
		h.logger.Debug(err.Error())
		return c.Next()
	}

	user := c.Locals("user").(models.User)

	if err := h.usecase.DeleteTodo(c.Context(), user.ID, uuid); err != nil {
		h.logger.Error(err.Error())
		return c.Next()
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "success",
	})
}
