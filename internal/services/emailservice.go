package services

import (
	"fmt"
	"log/slog"
)

type EmailService struct {
	sender EmailSender
	logger *slog.Logger
}

func NewEmailService(sender EmailSender, logger *slog.Logger) *EmailService {
	return &EmailService{
		sender: sender,
		logger: logger,
	}
}

func (es *EmailService) SendVerificationEmail(to string, code string) error {
	subject := "Verify your account"
	body := fmt.Sprintf("Your code is: %s", code)
	return es.sender.Send(to, subject, body)
}

func (es *EmailService) SendPasswordResetEmail(to string, code string) error {
	subject := "Password reset code"
	body := fmt.Sprintf("Your code is: %s", code)
	return es.sender.Send(to, subject, body)
}

func (es *EmailService) SendChangeEmailVerificationEmail(to string, code string) error {
	subject := "Verify if this is your new email"
	body := fmt.Sprintf("Your code is: %s", code)
	return es.sender.Send(to, subject, body)
}

func (es *EmailService) SendChangeEmailNotificationEmail(to string) error {
	subject := "Your email has been changed"
	body := "Your email has been changed"
	return es.sender.Send(to, subject, body)
}
