package services

import (
	"errors"
	"log/slog"
	"strings"

	"github.com/go-resty/resty/v2"
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

type MailtrapEmailSender struct {
	logger     *slog.Logger
	domain     string
	token      string
	restClient *resty.Client
}

func NewMailtrapEmailSender(logger *slog.Logger, domain, token string) *MailtrapEmailSender {
	client := resty.New()
	return &MailtrapEmailSender{
		logger:     logger,
		domain:     domain,
		token:      token,
		restClient: client,
	}
}

var ErrEmailSendingDisabled = errors.New("email sending is temporarily disabled")

func (mt *MailtrapEmailSender) Send(from, to, subject, body string) error {
	// temporarily disable email sending by returning an error here due to bots
	// until I implement reCAPTCHA
	return ErrEmailSendingDisabled

	/*
		requestBody := map[string]any{
			"from": map[string]string{
				"email": from,
			},
			"to": []map[string]string{
				{
					"email": to,
				},
			},
			"subject": subject,
			"text":    body,
			"html":    body,
		}

		resp, err := mt.restClient.R().
			SetBody(requestBody).
			SetHeader("Accept", "application/json").
			SetAuthToken(mt.token).
			Post("https://send.api.mailtrap.io/api/send")

		if err != nil {
			mt.logger.Error("An error occured while sending email with Mailtrap", "error", err)
			return err
		}

		if resp.IsError() {
			mt.logger.Error("An error occured while sending email with Mailtrap",
				"status", resp.Status(),
				"body", string(resp.Body()),
			)
			return fmt.Errorf("mailtrap returned %s", resp.Status())
		}

		return nil
	*/
}
