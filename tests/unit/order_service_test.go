package unit

import (
	"context"
	"testing"

	"ecommerce-service/internal/config"
	"ecommerce-service/internal/models"
	"ecommerce-service/internal/repositories"
	"ecommerce-service/internal/services"
	"ecommerce-service/tests"

	"github.com/stretchr/testify/assert"
)

func TestOrderService(t *testing.T) {
	//mockNotificationService := new(MockNotificationService)
	db := tests.SetupTestDB(t)

	// Setup dependencies
	orderRepo := repositories.NewOrderRepository(db)
	productRepo := repositories.NewProductRepository(db)
	customerRepo := repositories.NewCustomerRepository(db)
	cfg := config.NotificationConfig{}
	service := services.NewOrderService(orderRepo, *productRepo, *customerRepo, *services.NewNotificationService(cfg))

	t.Run("Create Order", func(t *testing.T) {
		order := &models.Order{
			CustomerID: "test-customer",
			Items: []models.OrderItem{
				{
					ProductID: "test-product",
					Quantity:  2,
					Price:     10.99,
				},
			},
		}

		err := service.CreateOrder(context.Background(), order)
		assert.NoError(t, err)
		assert.Equal(t, 21.98, order.TotalAmount) // 10.99 * 2
	})
}
