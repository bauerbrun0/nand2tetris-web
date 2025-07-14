package services

import (
	"errors"
	"log/slog"
	"strings"
)

type EmailSender interface {
	Send(to, subject, body string) error
}

type ConsoleEmailSender struct {
	logger *slog.Logger
}

func NewConsoleEmailSender(logger *slog.Logger) *ConsoleEmailSender {
	return &ConsoleEmailSender{
		logger: logger,
	}
}

func (es *ConsoleEmailSender) Send(to, subject, body string) error {
	if strings.Contains(body, "error:") {
		return errors.New("Error while sending email")
	}

	es.logger.Info("Sent email", "to", to, "subject", subject, "body", body)
	return nil
}
