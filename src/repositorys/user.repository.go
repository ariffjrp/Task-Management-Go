package repositorys

import (
	"context"
	"task_management/src/entity"

	"gorm.io/gorm"
)

type UserRepositorys interface {
	Register(ctx context.Context, user entity.User, account entity.Account) (entity.User, error)
}

type userRepositorys struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *userRepositorys {
	return &userRepositorys{db}
}

func (r *userRepositorys) Register(ctx context.Context, user entity.User, account entity.Account) (entity.User, error) {
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&user).Error; err != nil {
			return err
		}

		account.UserId = user.Id

		if err := tx.Create(&account).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return entity.User{}, err
	}

	r.db.Model(&user).Association("Account").Find(&user.Account)

	return user, nil
}
