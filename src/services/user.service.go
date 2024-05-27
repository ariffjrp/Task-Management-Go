package services

import (
	"context"
	"errors"
	"time"

	"task_management/src/entity"
	"task_management/src/repositorys"

	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Register(ctx context.Context, user *entity.User, account *entity.Account) (entity.User, error)
}

type userService struct {
	userRepository repositorys.UserRepositorys
}

// NewUserService adalah konstruktor untuk UserService
func NewUserService(userRepo repositorys.UserRepositorys) UserService {
	return &userService{
		userRepository: userRepo,
	}
}

// Register mendaftarkan pengguna baru
func (s *userService) Register(ctx context.Context, user *entity.User, account *entity.Account) (entity.User, error) {
	// Validasi input
	if user.Email == "" || user.Password == "" {
		return entity.User{}, errors.New("email and password are required")
	}

	// Hash kata sandi
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return entity.User{}, err
	}
	user.Password = string(hashedPassword)

	// Isi informasi tambahan
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	account.UserId = user.Id
	account.CreatedAt = now
	account.UpdatedAt = now

	// Simpan pengguna dan akun ke database
	registeredUser, err := s.userRepository.Register(ctx, *user, *account)
	if err != nil {
		return entity.User{}, err
	}

	return registeredUser, nil
}
