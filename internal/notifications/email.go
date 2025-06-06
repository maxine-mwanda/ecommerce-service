// internal/notifications/email.go
package notifications

import (
	"context"
	"ecommerce-service/internal/config"
	"fmt"
	"log"
	"net/smtp"
)

type EmailSender interface {
	Send(ctx context.Context, to, subject, body string) error
}

type SMTPEmailSender struct {
	cfg config.EmailConfig
}

func NewSMTPEmailSender(cfg config.EmailConfig) *SMTPEmailSender {
	return &SMTPEmailSender{cfg: cfg}
}

func (s *SMTPEmailSender) Send(ctx context.Context, to, subject, body string) error {
	auth := smtp.PlainAuth(
		"",
		s.cfg.SMTPUser,
		s.cfg.SMTPPassword,
		s.cfg.SMTPHost,
	)

	msg := []byte(fmt.Sprintf("To: %s\r\nSubject: %s\r\n\r\n%s\r\n", to, subject, body))

	addr := fmt.Sprintf("%s:%d", s.cfg.SMTPHost, s.cfg.SMTPPort)

	err := smtp.SendMail(
		addr,
		auth,
		s.cfg.FromEmail,
		[]string{to},
		msg,
	)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	log.Printf("✅ Email sent to %s | Subject: %s", to, subject)
	return nil
}
