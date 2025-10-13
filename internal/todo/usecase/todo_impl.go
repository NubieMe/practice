package usecase

import (
	"context"
	"errors"
	"practice/internal/todo/repository"
	"practice/models"
	"practice/pkg/logger"
	"practice/pkg/pagination"
	"practice/pkg/validator"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrNotFound  = errors.New("todo not found")
	ErrForbidden = errors.New("Forbidden")
)

type TodoUsecaseImpl struct {
	repo      repository.TodoRepo
	validator *validator.CustomValidator
	logger    *logger.Logger
}

func NewTodoUsecase(repo repository.TodoRepo, validator *validator.CustomValidator, logger *logger.Logger) TodoUsecase {
	return &TodoUsecaseImpl{
		repo:      repo,
		validator: validator,
		logger:    logger,
	}
}

func (u *TodoUsecaseImpl) AddTodo(ctx context.Context, todo *models.TodoRequest) error {
	err := u.validator.Validate(todo)
	if err != nil {
		u.logger.Debug(err.Error())
		return err
	}

	todoModel := &models.Todo{
		Title:  todo.Title,
		Todo:   todo.Todo,
		Check:  todo.Check,
		Images: todo.Images,
		UserID: todo.UserID,
	}

	return u.repo.AddTodo(ctx, todoModel)
}

func (u *TodoUsecaseImpl) UpdateTodo(ctx context.Context, todo *models.Todo) error {
	err := u.validator.Validate(todo)
	if err != nil {
		u.logger.Debug(err.Error())
		return err
	}

	existing, err := u.repo.GetTodo(ctx, todo.ID)
	if err != nil {
		u.logger.Debug(err.Error())
		return err
	}

	todo.UserID = existing.UserID
	todo.CreatedAt = existing.CreatedAt
	// todo.CreatedBy = existing.CreatedBy

	return u.repo.UpdateTodo(ctx, todo)
}

func (u *TodoUsecaseImpl) GetTodo(ctx context.Context, uuid uuid.UUID) (*models.Todo, error) {
	existing, err := u.repo.GetTodo(ctx, uuid)

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			u.logger.Debug(err.Error())
			return nil, ErrNotFound
		}

		u.logger.Debug(err.Error())
		return nil, err
	}

	return existing, nil
}

func (u *TodoUsecaseImpl) GetTodos(
	ctx context.Context,
	params *pagination.PaginationParams,
	userID uuid.UUID,
	key string,
	value ...interface{},
) ([]*models.Todo, error) {
	p := pagination.NewPagination(params)

	return u.repo.GetTodos(ctx, p, &userID, key, value...)
}

func (u *TodoUsecaseImpl) DeleteTodo(ctx context.Context, userID uuid.UUID, uuid uuid.UUID) error {
	existing, err := u.repo.GetTodo(ctx, uuid)
	if err != nil {
		u.logger.Debug(err.Error())
		return err
	}

	if existing.UserID != userID {
		u.logger.Debug("Forbidden when delete")
		return ErrForbidden
	}

	if existing == nil {
		u.logger.Debug("todo not found when delete")
		return ErrNotFound
	}

	return u.repo.DeleteTodo(ctx, uuid)
}
