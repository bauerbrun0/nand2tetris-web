package services

import (
	"bytes"
	"context"
	"log/slog"
	"time"

	"github.com/a-h/templ"
	"github.com/bauerbrun0/nand2tetris-web/ui/emails"
	"github.com/bauerbrun0/nand2tetris-web/ui/emails/emailchangeemail"
	"github.com/bauerbrun0/nand2tetris-web/ui/emails/emailchangenotificationemail"
	"github.com/bauerbrun0/nand2tetris-web/ui/emails/passwordresetemail"
	"github.com/bauerbrun0/nand2tetris-web/ui/emails/verificationemail"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type EmailService interface {
	SendVerificationEmail(to string, code string) error
	SendPasswordResetEmail(to string, code string) error
	SendChangeEmailVerificationEmail(to string, code string) error
	SendChangeEmailNotificationEmail(to string) error
}

type emailService struct {
	sender       EmailSender
	logger       *slog.Logger
	noreplyEmail string
	localizer    *i18n.Localizer
	baseUrl      string
}

func NewEmailService(sender EmailSender, logger *slog.Logger, localizer *i18n.Localizer, noreplyEmail string, baseUrl string) EmailService {
	return &emailService{
		sender:       sender,
		logger:       logger,
		noreplyEmail: noreplyEmail,
		localizer:    localizer,
		baseUrl:      baseUrl,
	}
}

func (es *emailService) newEmailData() emails.EmailData {
	return emails.EmailData{
		CurrentYear: time.Now().Year(),
		Localizer:   es.localizer,
		BaseUrl:     es.baseUrl,
	}
}

func renderToString(ctx context.Context, c templ.Component) (string, error) {
	var buf bytes.Buffer
	err := c.Render(ctx, &buf)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (es *emailService) SendVerificationEmail(to string, code string) error {
	baseEmailData := es.newEmailData()
	email := verificationemail.Email(verificationemail.VerificationEmailData{
		Code:      code,
		EmailData: baseEmailData,
	})
	body, err := renderToString(context.Background(), email)
	if err != nil {
		return err
	}
	subject := baseEmailData.T("verification_email.subject")
	return es.sender.Send(es.noreplyEmail, to, subject, body)
}

func (es *emailService) SendPasswordResetEmail(to string, code string) error {
	baseEmailData := es.newEmailData()
	email := passwordresetemail.Email(passwordresetemail.PasswordResetEmailData{
		Code:      code,
		EmailData: baseEmailData,
	})
	body, err := renderToString(context.Background(), email)
	if err != nil {
		return err
	}
	subject := baseEmailData.T("password_reset_email.subject")
	return es.sender.Send(es.noreplyEmail, to, subject, body)
}

func (es *emailService) SendChangeEmailVerificationEmail(to string, code string) error {
	baseEmailData := es.newEmailData()
	email := emailchangeemail.Email(emailchangeemail.EmailChangeEmailData{
		Code:      code,
		EmailData: baseEmailData,
	})
	body, err := renderToString(context.Background(), email)
	if err != nil {
		return err
	}
	subject := baseEmailData.T("email_change_email.subject")
	return es.sender.Send(es.noreplyEmail, to, subject, body)
}

func (es *emailService) SendChangeEmailNotificationEmail(to string) error {
	baseEmailData := es.newEmailData()
	email := emailchangenotificationemail.Email(baseEmailData)
	body, err := renderToString(context.Background(), email)
	if err != nil {
		return err
	}
	subject := baseEmailData.T("email_change_notification_email.subject")
	return es.sender.Send(es.noreplyEmail, to, subject, body)
}
