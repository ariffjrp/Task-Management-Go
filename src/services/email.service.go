package services

import (
	"fmt"
	"net/smtp"
	configs "task_management/src/configs"
	"task_management/src/utils"
)

type EmailService struct {
	EmailConfig configs.EmailConfig
	TOTPService *TOTPService
}

func NewEmailService(emailConfig configs.EmailConfig, totpService *TOTPService) *EmailService {
	return &EmailService{
		EmailConfig: emailConfig,
		TOTPService: totpService,
	}
}

func (s *EmailService) SendEmail(to, subject, body string) error {
	from := s.EmailConfig.SMPTUsername
	password := s.EmailConfig.SMPTPassword
	smtpHost := s.EmailConfig.SMTPHost
	smtpPort := s.EmailConfig.SMTPPort

	auth := smtp.PlainAuth("", from, password, smtpHost)

	message := []byte(fmt.Sprintf("Subject: %s\n\n%s", subject, body))

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, message)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

func (s *EmailService) SendVerificationEmail(to, code string) error {
	subject := "Email Verification"
	body := utils.GenerateOTPEmailBody(code)
	return s.SendEmail(to, subject, body)
}

func (s *EmailService) SendForgetPasswordEmail(to, link string) error {
	subject := "Reset Password"
	body := utils.GenerateForgetPasswordEmailBody(link)
	return s.SendEmail(to, subject, body)
}
