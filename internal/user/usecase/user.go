package usecase

import (
	"context"
	"practice/models"
	"practice/pkg/pagination"
)

type UserUsecase interface {
	Register(ctx context.Context, user *models.UserRegister) error
	Login(ctx context.Context, user *models.UserLogin) (string, error)
	GetUser(ctx context.Context, uuidStr string) (*models.User, error)
	GetUsers(ctx context.Context, params *pagination.PaginationParams, key string, value ...interface{}) ([]*models.User, error)
}
