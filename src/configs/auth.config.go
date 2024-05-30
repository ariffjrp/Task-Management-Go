package configs

import (
	"fmt"
	"time"

	"github.com/joho/godotenv"
)

type JWTConfig struct {
	Secret               string
	JWTExpiration        time.Duration
	JWTRefreshExpiration time.Duration
}

func LoadConfigJWT() *JWTConfig {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env JWT")
	}

	return &JWTConfig{
		Secret:               getEnv("SECRET_KEY"),
		JWTExpiration:        time.Hour * 1,
		JWTRefreshExpiration: time.Hour * 24 * 7,
	}
}
