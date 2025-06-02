// internal/notifications/email.go
package notifications

import (
	"context"
	"ecommerce-service/internal/config"
	"log"
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
	// In a real implementation, this would send an email via SMTP
	log.Printf("Sending email to %s: %s - %s\n", to, subject, body)
	return nil
}
