package handler

import "github.com/gofiber/fiber/v2"

type TodoHandler interface {
	GetTodos(c *fiber.Ctx) error
	GetTodo(c *fiber.Ctx) error
	AddTodo(c *fiber.Ctx) error
	UpdateTodo(c *fiber.Ctx) error
	DeleteTodo(c *fiber.Ctx) error
}
