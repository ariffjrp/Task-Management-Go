package services

import (
	"context"
	"errors"
	"time"

	"task_management/src/entity"
	repository "task_management/src/repositorys"

	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Register(ctx context.Context, user *entity.User, account *entity.Account) (entity.User, entity.Account, error)
}

type userService struct {
	userRepository repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{
		userRepository: userRepo,
	}
}

func (s *userService) Register(ctx context.Context, user *entity.User, account *entity.Account) (entity.User, entity.Account, error) {
	if user.Email == "" || user.Password == "" {
		return entity.User{}, entity.Account{}, errors.New("email and password are required")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return entity.User{}, entity.Account{}, err
	}
	user.Password = string(hashedPassword)

	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	account.CreatedAt = now
	account.UpdatedAt = now

	registeredUser, registeredAccount, err := s.userRepository.Register(ctx, *user, *account)
	if err != nil {
		return entity.User{}, entity.Account{}, err
	}

	return registeredUser, registeredAccount, nil
}
