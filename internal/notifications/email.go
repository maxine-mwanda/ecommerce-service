// internal/notifications/email.go
package notifications

import (
	"context"
	"fmt"
	"log"
	"net/smtp"

	"ecommerce-service/internal/config"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
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
	from := s.cfg.SMTPEmail

	// OAuth2 config setup
	config := &oauth2.Config{
		ClientID:     s.cfg.SMTPClientID,
		ClientSecret: s.cfg.SMTPClientSecret,
		Endpoint:     google.Endpoint,
		Scopes:       []string{"https://mail.google.com/"},
		RedirectURL:  "urn:ietf:wg:oauth:2.0:oob",
	}

	token := &oauth2.Token{
		RefreshToken: s.cfg.SMTPRefreshToken,
	}

	// Retrieve access token
	tokenSource := config.TokenSource(ctx, token)
	newToken, err := tokenSource.Token()
	if err != nil {
		return fmt.Errorf("failed to retrieve access token: %w", err)
	}

	accessToken := newToken.AccessToken

	// Prepare auth with access token
	auth := smtp.PlainAuth("", from, accessToken, "smtp.gmail.com")

	// Email message
	msg := []byte("To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"MIME-version: 1.0;\r\nContent-Type: text/plain; charset=\"UTF-8\";\r\n\r\n" +
		body + "\r\n")

	// Send the email
	err = smtp.SendMail("smtp.gmail.com:587", auth, from, []string{to}, msg)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	log.Printf("Email sent to %s successfully.\n", to)
	return nil
}
