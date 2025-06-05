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

/*func NewNotificationService(cfg config.NotificationConfig) *NotificationService {
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
*/
package services_test

import (
	"context"
	"testing"

	"ecommerce-service/internal/models"
	"ecommerce-service/internal/repositories"
	"ecommerce-service/internal/services"
	"ecommerce-service/tests"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestProductService_CreateProduct(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock: %v", err)
	}
	defer db.Close()

	productRepo := repositories.NewProductRepository(db)
	categoryRepo := repositories.NewCategoryRepository(db)
	service := services.NewProductService(productRepo, categoryRepo)

	t.Run("Success", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO products").
			WillReturnResult(sqlmock.NewResult(1, 1))

		product := &models.Product{
			Name:       "Test Product",
			Price:      10.99,
			CategoryID: "cat-123",
		}

		err := service.CreateProduct(context.Background(), product)
		assert.NoError(t, err)
		assert.NotEmpty(t, product.ID)
	})

	// Add more test cases...
}