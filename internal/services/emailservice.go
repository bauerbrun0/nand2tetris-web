package services

import "fmt"

type EmailService struct {
	sender EmailSender
}

func NewEmailService(sender EmailSender) *EmailService {
	return &EmailService{
		sender: sender,
	}
}

func (es *EmailService) SendVerificationEmail(to string, code string) error {
	subject := "Verify your account"
	body := fmt.Sprintf("Your code is: %s", code)
	return es.sender.Send(to, subject, body)
}
