package usecase

import (
	"context"
	"fmt"
	"practice/env"
	"practice/internal/user/repository"
	"practice/models"
	"practice/pkg/logger"
	"practice/pkg/pagination"
	"practice/pkg/validator"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserUsecaseImpl struct {
	repo      repository.UserRepo
	validator *validator.CustomValidator
	logger    *logger.Logger
}

func NewUserUsecase(repo repository.UserRepo, validator *validator.CustomValidator, logger *logger.Logger) UserUsecase {
	return &UserUsecaseImpl{
		repo:      repo,
		validator: validator,
		logger:    logger,
	}
}

func (u *UserUsecaseImpl) Register(ctx context.Context, user *models.UserRegister) error {
	err := u.validator.Validate(user)
	if err != nil {
		u.logger.Debug(err.Error())
		return err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		u.logger.Debug(err.Error())
		return err
	}
	user.Password = string(hash)

	userModel := &models.User{
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
	}

	if err := u.repo.AddUser(ctx, userModel); err != nil {
		u.logger.Debug(err.Error())
		return err
	}

	return nil
}

func (u *UserUsecaseImpl) Login(ctx context.Context, user *models.UserLogin) (string, error) {
	err := u.validator.Validate(user)
	if err != nil {
		u.logger.Debug(err.Error())
		return "", err
	}

	fmt.Println("user", user)
	existing, err := u.repo.GetUserByEmail(ctx, user.Email)
	if err != nil {
		u.logger.Debug(err.Error())
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(existing.Password), []byte(user.Password)); err != nil {
		u.logger.Debug(err.Error())
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    existing.ID,
		"email": existing.Email,
		"name":  existing.Name,
	})
	signed, err := token.SignedString([]byte(env.JWTSecretKey))

	if err != nil {
		u.logger.Debug(err.Error())
		return "", err
	}

	return signed, nil
}

func (u *UserUsecaseImpl) GetUser(ctx context.Context, uuidStr string) (*models.User, error) {
	uuid, err := uuid.Parse(uuidStr)
	if err != nil {
		u.logger.Debug(err.Error())
		return nil, err
	}

	return u.repo.GetUser(ctx, uuid)
}

func (u *UserUsecaseImpl) GetUsers(ctx context.Context, params *pagination.PaginationParams, key string, value ...interface{}) ([]*models.User, error) {
	p := pagination.NewPagination(params)

	return u.repo.GetUsers(ctx, p, key, value...)
}
