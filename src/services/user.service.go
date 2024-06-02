package services

import (
	"context"
	"errors"
	"time"

	"task_management/src/configs"
	"task_management/src/entity"
	repository "task_management/src/repositorys"
	"task_management/src/utils"

	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Register(ctx context.Context, user *entity.User, account *entity.Account) (entity.User, entity.Account, string, string, error)
	Login(ctx context.Context, loginRequest *entity.UserLogin) (string, string, error)
	RefreshToken(ctx context.Context, refreshToken string) (string, string, error)
}

type userService struct {
	userRepository repository.UserRepository
	jwtConfig      *configs.JWTConfig
}

func NewUserService(userRepo repository.UserRepository, jwtConfig *configs.JWTConfig) UserService {
	return &userService{
		userRepository: userRepo,
		jwtConfig:      jwtConfig,
	}
}

func (s *userService) Register(ctx context.Context, user *entity.User, account *entity.Account) (entity.User, entity.Account, string, string, error) {
	if user.Email == "" || user.Password == "" {
		return entity.User{}, entity.Account{}, "", "", errors.New("email and password are required")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return entity.User{}, entity.Account{}, "", "", err
	}
	user.Password = string(hashedPassword)

	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	account.CreatedAt = now
	account.UpdatedAt = now

	registeredUser, registeredAccount, err := s.userRepository.Register(ctx, *user, *account)
	if err != nil {
		return entity.User{}, entity.Account{}, "", "", err
	}

	accessToken, err := utils.GenerateToken(&registeredUser, s.jwtConfig)
	if err != nil {
		return entity.User{}, entity.Account{}, "", "", err
	}

	refreshToken, err := utils.GenerateRefreshToken(&registeredUser, s.jwtConfig)
	if err != nil {
		return entity.User{}, entity.Account{}, "", "", err
	}

	refreshTokenEntity := entity.RefreshToken{
		UserId:     registeredUser.Id,
		Token:      refreshToken,
		ExpiryDate: time.Now().Add(s.jwtConfig.JWTRefreshExpiration),
	}

	err = s.userRepository.SaveRefreshToken(ctx, refreshTokenEntity)
	if err != nil {
		return entity.User{}, entity.Account{}, "", "", err
	}

	return registeredUser, registeredAccount, accessToken, refreshToken, nil
}

func (s *userService) Login(ctx context.Context, loginRequest *entity.UserLogin) (string, string, error) {
	dbUser, err := s.userRepository.Login(ctx, loginRequest.Email)
	if err != nil {
		return "", "", err
	}

	if dbUser.Email == "" {
		return "", "", errors.New("user not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(loginRequest.Password))
	if err != nil {
		return "", "", errors.New("wrong email or password")
	}

	accessToken, err := utils.GenerateToken(&dbUser, s.jwtConfig)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := utils.GenerateRefreshToken(&dbUser, s.jwtConfig)
	if err != nil {
		return "", "", err
	}

	refreshTokenEntity := entity.RefreshToken{
		UserId:     dbUser.Id,
		Token:      refreshToken,
		ExpiryDate: time.Now().Add(s.jwtConfig.JWTRefreshExpiration),
	}

	err = s.userRepository.SaveRefreshToken(ctx, refreshTokenEntity)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (s *userService) RefreshToken(ctx context.Context, refreshToken string) (string, string, error) {
	user, err := s.userRepository.GetUserByRefreshToken(ctx, refreshToken)
	if err != nil {
		return "", "", errors.New("invalid refresh token")
	}

	newAccessToken, err := utils.GenerateToken(&user, s.jwtConfig)
	if err != nil {
		return "", "", err
	}

	newRefreshToken, err := utils.GenerateRefreshToken(&user, s.jwtConfig)
	if err != nil {
		return "", "", err
	}

	err = s.userRepository.DeleteRefreshToken(ctx, refreshToken)
	if err != nil {
		return "", "", err
	}

	newRefreshTokenEntity := entity.RefreshToken{
		UserId:     user.Id,
		Token:      newRefreshToken,
		ExpiryDate: time.Now().Add(s.jwtConfig.JWTRefreshExpiration),
	}

	err = s.userRepository.SaveRefreshToken(ctx, newRefreshTokenEntity)
	if err != nil {
		return "", "", err
	}

	return newAccessToken, newRefreshToken, nil
}
