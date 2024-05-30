package utils

import (
	"task_management/src/configs"
	"task_management/src/entity"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func GenerateToken(user *entity.User, config *configs.JWTConfig) (string, error) {
	claims := &entity.Claims{
		UserId: user.Id.String(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(config.JWTExpiration)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.Secret))
}

func GenerateRefreshToken(user *entity.User, config *configs.JWTConfig) (string, error) {
	claims := &entity.Claims{
		UserId: user.Id.String(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(config.JWTRefreshExpiration)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.Secret))
}
