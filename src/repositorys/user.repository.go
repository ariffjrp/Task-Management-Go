package repository

import (
	"context"
	"task_management/src/entity"

	"gorm.io/gorm"
)

type UserRepository interface {
	Register(ctx context.Context, user entity.User, account entity.Account) (entity.User, entity.Account, error)
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
