package repository

import (
	"context"
	"task_management/src/entity"

	"gorm.io/gorm"
)

type AuthRepository interface {
	SaveOAuth2User(ctx context.Context, user *entity.User) error
	GetUserByOAuth2ID(ctx context.Context, oauth2ID string) (*entity.User, error)
}

type authRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) AuthRepository {
	return &authRepository{db}
}

func (r *authRepository) SaveOAuth2User(ctx context.Context, user *entity.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

func (r *authRepository) GetUserByOAuth2ID(ctx context.Context, oauth2ID string) (*entity.User, error) {
	var user entity.User
	err := r.db.WithContext(ctx).Where("o_auth2_id = ?", oauth2ID).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
