package services

import (
	"crypto/rand"
	"encoding/base32"
	"errors"
	"time"

	"github.com/pquerna/otp/totp"
)

type TOTPService struct{}

func NewTOTPService() *TOTPService {
	return &TOTPService{}
}

func (s *TOTPService) GenerateSecret() (string, error) {
	secretBytes := make([]byte, 10)
	_, err := rand.Read(secretBytes)
	if err != nil {
		return "", err
	}

	secret := base32.StdEncoding.EncodeToString(secretBytes)

	return secret, nil
}

func (s *TOTPService) GenerateOTP(secret string) (string, error) {
	totpCode, err := totp.GenerateCode(secret, time.Now())
	if err != nil {
		return "", err
	}

	return totpCode, nil
}

func (s *TOTPService) VerifyOTP(secret, otp string) error {
	if !totp.Validate(otp, secret) {
		return errors.New("invalid OTP")
	}
	return nil
}
