package handler

import (
	"practice/internal/user/usecase"
	"practice/models"
	"practice/pkg/bus"
	"practice/pkg/logger"
	"practice/pkg/pagination"

	"github.com/gofiber/fiber/v2"
)

type UserHandlerImpl struct {
	usecase usecase.UserUsecase
	logger  *logger.Logger
	event   *bus.EventBus
}

func NewUserHandler(usecase usecase.UserUsecase, logger *logger.Logger, event *bus.EventBus) UserHandler {
	return &UserHandlerImpl{
		usecase: usecase,
		logger:  logger,
		event:   event,
	}
}

func (h *UserHandlerImpl) Handle(event bus.Event) {
	h.logger.Info("Todo event: %s - %v", event.Type, event.Payload)
}

func (h *UserHandlerImpl) Register(c *fiber.Ctx) error {
	request := new(models.UserRegister)
	if err := c.BodyParser(request); err != nil {
		return c.Next()
	}

	if err := h.usecase.Register(c.Context(), request); err != nil {
		return c.Next()
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "success",
	})
}

func (h *UserHandlerImpl) Login(c *fiber.Ctx) error {
	request := new(models.UserLogin)
	if err := c.BodyParser(request); err != nil {
		h.logger.Debug(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid request",
		})
	}

	token, err := h.usecase.Login(c.Context(), request)
	if err != nil {
		h.logger.Error(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "email/password invalid",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "success",
		"token":   token,
	})
}

func (h *UserHandlerImpl) GetUser(c *fiber.Ctx) error {
	uuidStr := c.Params("id")

	user, err := h.usecase.GetUser(c.Context(), uuidStr)
	if err != nil {
		return c.Next()
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "success",
		"data":    user,
	})
}

func (h *UserHandlerImpl) GetUsers(c *fiber.Ctx) error {
	name := c.Query("name", "")

	var pagParams models.PaginationRequest
	if err := c.QueryParser(&pagParams); err != nil {
		h.logger.Debug(err.Error())
		return c.Next()
	}

	params := &pagination.PaginationParams{
		Page:  pagParams.Page,
		Limit: pagParams.Limit,
		Sort:  pagParams.Sort,
	}

	users, err := h.usecase.GetUsers(c.Context(), params, "name", name)
	if err != nil {
		return c.Next()
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "success",
		"data":    users,
	})
}
