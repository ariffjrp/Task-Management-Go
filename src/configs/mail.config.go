package configs

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type EmailConfig struct {
	SMTPHost     string
	SMTPPort     string
	SMPTUsername string
	SMPTPassword string
}

func LoadConfigEmail() EmailConfig {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env Mail")
	}

	return EmailConfig{
		SMTPHost:     getEnv("SMPT_HOST"),
		SMTPPort:     getEnv("SMPT_PORT"),
		SMPTUsername: getEnv("SMPT_USERNAME"),
		SMPTPassword: getEnv("SMPT_PASSWORD"),
	}
}

func getEnv(key string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return "Error get Key error"
	}
	return value
}
