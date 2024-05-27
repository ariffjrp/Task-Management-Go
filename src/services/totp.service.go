package services

import (
	"time"

	"github.com/pquerna/otp/totp"
)

type TOTPService struct{}

func NewTOTPService() *TOTPService {
	return &TOTPService{}
}

func (s *TOTPService) Generate(secret string) (string, error) {
	totpCode, err := totp.GenerateCode(secret, time.Now())
	if err != nil {
		return "", err
	}

	return totpCode, nil
}

func (s *TOTPService) VerifyOTP(secret, otp string) bool {
	return totp.Validate(otp, secret)
}
