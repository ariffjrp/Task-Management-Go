package services

import (
	"context"
	"encoding/json"
	"fmt"
	"task_management/src/configs"
	"task_management/src/entity"
	repository "task_management/src/repositorys"

	"golang.org/x/oauth2"
	"gorm.io/gorm"
)

type AuthService interface {
	GetGoogleOAuth2Config() *oauth2.Config
	HandleGoogleCallback(ctx context.Context, state string, code string) (*entity.User, error)
	GetJWTConfig() *configs.JWTConfig
}

type authService struct {
	repository repository.AuthRepository
	config     *configs.JWTConfig
}

func NewAuthService(repository repository.AuthRepository, config *configs.JWTConfig) AuthService {
	return &authService{repository: repository, config: config}
}

func (s *authService) GetGoogleOAuth2Config() *oauth2.Config {
	return configs.GoogleOAuth2Config
}

func (s *authService) HandleGoogleCallback(ctx context.Context, state string, code string) (*entity.User, error) {
	token, err := configs.GoogleOAuth2Config.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange token: %w", err)
	}

	client := configs.GoogleOAuth2Config.Client(ctx, token)
	userInfoResp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer userInfoResp.Body.Close()

	var userInfo struct {
		Id    string `json:"id"`
		Email string `json:"email"`
	}

	if err := json.NewDecoder(userInfoResp.Body).Decode(&userInfo); err != nil {
		return nil, fmt.Errorf("failed to decode user info: %w", err)
	}

	user, err := s.repository.GetUserByOAuth2ID(ctx, userInfo.Id)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("failed to get user by OAuth2 ID: %w", err)
	}

	if user == nil {
		user = &entity.User{
			OAuth2Id: userInfo.Id,
			Email:    userInfo.Email,
			Provider: "google",
		}
		if err := s.repository.SaveOAuth2User(ctx, user); err != nil {
			return nil, fmt.Errorf("failed to save user: %w", err)
		}
	}

	return user, nil
}

func (s *authService) GetJWTConfig() *configs.JWTConfig {
	return s.config
}
