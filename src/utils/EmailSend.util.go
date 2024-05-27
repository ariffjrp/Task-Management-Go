package utils

import "fmt"

func GenerateOTPEmailBody(code string) string {
	return fmt.Sprintf("Your verification code is: %s", code)
}

func GenerateForgetPasswordEmailBody(link string) string {
	return fmt.Sprintf("Click the following link to reset your password: %s", link)
}
