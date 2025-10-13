package usecase

import (
	"context"
	"practice/models"
	"practice/pkg/pagination"

	"github.com/google/uuid"
)

type TodoUsecase interface {
	AddTodo(ctx context.Context, todo *models.TodoRequest) error
	UpdateTodo(ctx context.Context, todo *models.Todo) error
	GetTodo(ctx context.Context, uuid uuid.UUID) (*models.Todo, error)
	GetTodos(ctx context.Context, params *pagination.PaginationParams, userID uuid.UUID, key string, value ...interface{}) ([]*models.Todo, error)
	DeleteTodo(ctx context.Context, userID uuid.UUID, uuid uuid.UUID) error
}
