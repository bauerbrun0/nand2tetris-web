package services

import (
	"context"
	"errors"
	"log/slog"
	"strings"
	"time"

	"github.com/mailgun/mailgun-go/v5"
)

type EmailSender interface {
	Send(from, to, subject, body string) error
}

type ConsoleEmailSender struct {
	logger *slog.Logger
}

func NewConsoleEmailSender(logger *slog.Logger) *ConsoleEmailSender {
	return &ConsoleEmailSender{
		logger: logger,
	}
}

func (es *ConsoleEmailSender) Send(from, to, subject, body string) error {
	if strings.Contains(body, "error:") {
		return errors.New("Error while sending email")
	}

	es.logger.Info("Sent email", "from", from, "to", to, "subject", subject, "body", body)
	return nil
}

type MailGunEmailSender struct {
	logger *slog.Logger
	domain string
	client *mailgun.Client
}

func NewMailGunEmailSender(logger *slog.Logger, domain, apiKey string) *MailGunEmailSender {
	mg := mailgun.NewMailgun(apiKey)
	return &MailGunEmailSender{
		logger: logger,
		domain: domain,
		client: mg,
	}
}

func (mg *MailGunEmailSender) Send(from, to, subject, body string) error {
	message := mailgun.NewMessage(mg.domain, from, subject, body, to)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	_, err := mg.client.Send(ctx, message)
	if err != nil {
		mg.logger.Error("An error occured while sending email", "error", err)
		return err
	}
	return nil
}
