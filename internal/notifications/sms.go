// internal/notifications/sms.go
package notifications

import (
	"context"
	"ecommerce-service/internal/config"
	"log"
)

type SMSSender interface {
	Send(ctx context.Context, phone, message string) error
}

type AfricaTalkingSMSSender struct {
	cfg config.SMSConfig
}

func NewAfricaTalkingSMSSender(cfg config.SMSConfig) *AfricaTalkingSMSSender {
	return &AfricaTalkingSMSSender{cfg: cfg}
}

func (s *AfricaTalkingSMSSender) Send(ctx context.Context, phone, message string) error {
	// In a real implementation, this would call Africa's Talking API
	// For now, we'll just log it
	log.Printf("Sending SMS to %s: %s (sandbox: %v)\n", phone, message, s.cfg.Sandbox)
	return nil
}
