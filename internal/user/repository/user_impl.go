package repository

import (
	"context"
	"practice/models"
	"practice/pkg/pagination"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepoImpl struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) UserRepo {
	return &UserRepoImpl{
		db: db,
	}
}

func (r *UserRepoImpl) AddUser(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Model(&models.User{}).Create(user).Error
}

func (r *UserRepoImpl) UpdateUser(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Model(&models.User{}).Save(user).Error
}

func (r *UserRepoImpl) GetUser(ctx context.Context, uuid uuid.UUID) (*models.User, error) {
	var user *models.User
	err := r.db.WithContext(ctx).Model(&models.User{}).First(&user, "id = ?", uuid).Error
	return user, err
}

func (r *UserRepoImpl) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user *models.User
	err := r.db.WithContext(ctx).Model(&models.User{}).First(&user, "email = ?", email).Error
	return user, err
}

func (r *UserRepoImpl) GetUsers(ctx context.Context, pagi *pagination.Pagination, key string, value ...interface{}) ([]*models.User, error) {
	var users []*models.User

	query := r.db.WithContext(ctx)
	if len(value) > 0 {
		query = query.Where("LOWER ("+key+") LIKE LOWER(?)", "%"+value[0].(string)+"%")
	}

	paginated, err := pagination.Paginate(&models.User{}, pagi, query)
	if err != nil {
		return nil, err
	}

	if err := paginated.Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}
