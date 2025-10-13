package repository

import (
	"context"
	"practice/models"
	"practice/pkg/pagination"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TodoRepoImpl struct {
	db *gorm.DB
}

func NewTodoRepo(db *gorm.DB) TodoRepo {
	return &TodoRepoImpl{
		db: db,
	}
}

func (r *TodoRepoImpl) AddTodo(ctx context.Context, todo *models.Todo) error {
	return r.db.WithContext(ctx).Model(&models.Todo{}).Create(todo).Error
}

func (r *TodoRepoImpl) UpdateTodo(ctx context.Context, todo *models.Todo) error {
	return r.db.WithContext(ctx).Model(&models.Todo{}).Save(todo).Error
}

func (r *TodoRepoImpl) GetTodo(ctx context.Context, uuid uuid.UUID) (*models.Todo, error) {
	var todo *models.Todo
	err := r.db.WithContext(ctx).Model(&models.Todo{}).First(&todo, "id = ?", uuid).Error
	return todo, err
}

func (r *TodoRepoImpl) GetTodos(ctx context.Context, params *pagination.Pagination, userID *uuid.UUID, key string, value ...interface{}) ([]*models.Todo, error) {
	var todos []*models.Todo

	query := r.db.WithContext(ctx)
	if userID != nil {
		query = query.Where("user_id = ?", userID)
	}

	if len(value) > 0 && value[0] != "" {
		searchTerm := "%" + value[0].(string) + "%"
		query = query.Where(`
			title ILIKE ? OR
			EXISTS (
			SELECT 1 FROM unnest(todo) AS item 
			WHERE item ILIKE ?
		)`, searchTerm, searchTerm)
	}

	paginated, err := pagination.Paginate(&models.Todo{}, params, query)
	if err != nil {
		return nil, err
	}

	if err := paginated.Preload("User").Find(&todos).Error; err != nil {
		return nil, err
	}

	return todos, nil
}

func (r *TodoRepoImpl) DeleteTodo(ctx context.Context, uuid uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&models.Todo{}).Where("id = ?", uuid).Delete(&models.Todo{}).Error
}
