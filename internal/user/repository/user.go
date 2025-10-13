package repository

import (
	"context"
	"practice/models"
	"practice/pkg/pagination"

	"github.com/google/uuid"
)

type UserRepo interface {
	AddUser(ctx context.Context, user *models.User) error
	UpdateUser(ctx context.Context, user *models.User) error
	GetUser(ctx context.Context, uuid uuid.UUID) (*models.User, error)
	GetUsers(ctx context.Context, params *pagination.Pagination, key string, value ...interface{}) ([]*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
}
