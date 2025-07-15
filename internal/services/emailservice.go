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
