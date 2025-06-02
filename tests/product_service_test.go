// internal/services/notification.go
package services

import (
	"ecommerce-service/internal/config"
	"log"
)

type NotificationService struct {
	smsConfig   config.SMSConfig
	emailConfig config.EmailConfig
}

func NewNotificationService(cfg config.NotificationConfig) *NotificationService {
	return &NotificationService{
		smsConfig:   cfg.SMS,
		emailConfig: cfg.Email,
	}
}

func (s *NotificationService) SendOrderConfirmationSMS(phone, message string) error {
	// In a real implementation, this would call Africa's Talking API
	log.Printf("Sending SMS to %s: %s\n", phone, message)
	return nil
}

func (s *NotificationService) SendOrderNotificationEmail(email, subject, message string) error {
	// In a real implementation, this would send an email
	log.Printf("Sending email to %s: %s - %s\n", email, subject, message)
	return nil
}
