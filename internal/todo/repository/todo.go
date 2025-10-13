package repository

import (
	"context"
	"practice/models"
	"practice/pkg/pagination"

	"github.com/google/uuid"
)

type TodoRepo interface {
	AddTodo(ctx context.Context, todo *models.Todo) error
	UpdateTodo(ctx context.Context, todo *models.Todo) error
	GetTodo(ctx context.Context, uuid uuid.UUID) (*models.Todo, error)
	GetTodos(ctx context.Context, params *pagination.Pagination, userID *uuid.UUID, key string, value ...interface{}) ([]*models.Todo, error)
	DeleteTodo(ctx context.Context, uuid uuid.UUID) error
}
