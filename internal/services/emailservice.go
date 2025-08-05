package services

import (
	"fmt"
	"log/slog"
)

type EmailService interface {
	SendVerificationEmail(to string, code string) error
	SendPasswordResetEmail(to string, code string) error
	SendChangeEmailVerificationEmail(to string, code string) error
	SendChangeEmailNotificationEmail(to string) error
}

type emailService struct {
	sender EmailSender
	logger *slog.Logger
}

func NewEmailService(sender EmailSender, logger *slog.Logger) EmailService {
	return &emailService{
		sender: sender,
		logger: logger,
	}
}

func (es *emailService) SendVerificationEmail(to string, code string) error {
	subject := "Verify your account"
	body := fmt.Sprintf("Your code is: %s", code)
	return es.sender.Send(to, subject, body)
}

func (es *emailService) SendPasswordResetEmail(to string, code string) error {
	subject := "Password reset code"
	body := fmt.Sprintf("Your code is: %s", code)
	return es.sender.Send(to, subject, body)
}

func (es *emailService) SendChangeEmailVerificationEmail(to string, code string) error {
	subject := "Verify if this is your new email"
	body := fmt.Sprintf("Your code is: %s", code)
	return es.sender.Send(to, subject, body)
}

func (es *emailService) SendChangeEmailNotificationEmail(to string) error {
	subject := "Your email has been changed"
	body := "Your email has been changed"
	return es.sender.Send(to, subject, body)
}
