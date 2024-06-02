package repository

import (
	"context"
	"errors"
	"task_management/src/entity"

	"gorm.io/gorm"
)

type UserRepository interface {
	Register(ctx context.Context, user entity.User, account entity.Account) (entity.User, entity.Account, error)
	Login(ctx context.Context, email string) (entity.User, error)
	SaveRefreshToken(ctx context.Context, refreshToken entity.RefreshToken) error
	GetUserByRefreshToken(ctx context.Context, refreshToken string) (entity.User, error)
	DeleteRefreshToken(ctx context.Context, refreshToken string) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db}
}

func (r *userRepository) Register(ctx context.Context, user entity.User, account entity.Account) (entity.User, entity.Account, error) {
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&user).Error; err != nil {
			return err
		}

		account.UserId = user.Id

		// Create account
		if err := tx.Create(&account).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return entity.User{}, entity.Account{}, err
	}

	if err := r.db.WithContext(ctx).Model(&user).Association("Account").Find(&user.Account); err != nil {
		return entity.User{}, entity.Account{}, err
	}

	return user, account, nil
}

func (r *userRepository) Login(ctx context.Context, email string) (entity.User, error) {
	var user entity.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.User{}, nil
		}
		return entity.User{}, err
	}

	return user, nil
}

func (r *userRepository) SaveRefreshToken(ctx context.Context, refreshToken entity.RefreshToken) error {
	return r.db.WithContext(ctx).Create(&refreshToken).Error
}

func (r *userRepository) GetUserByRefreshToken(ctx context.Context, refreshToken string) (entity.User, error) {
	var token entity.RefreshToken
	err := r.db.WithContext(ctx).Where("token = ?", refreshToken).First(&token).Error
	if err != nil {
		return entity.User{}, err
	}

	var user entity.User
	err = r.db.WithContext(ctx).Where("id = ?", token.UserId).First(&user).Error
	return user, err
}

func (r *userRepository) DeleteRefreshToken(ctx context.Context, refreshToken string) error {
	return r.db.WithContext(ctx).Where("token = ?", refreshToken).Delete(&entity.RefreshToken{}).Error
}
