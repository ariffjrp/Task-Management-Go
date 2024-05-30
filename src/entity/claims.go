package entity

import (
	"github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	UserId string `json:"user_id"`
	jwt.RegisteredClaims
}
