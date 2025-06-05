package unit

import (
	"context"
	"testing"

	"ecommerce-service/internal/models"
	"ecommerce-service/internal/repositories"
	"ecommerce-service/internal/services"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestOrderService(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock: %v", err)
	}
	defer db.Close()

	orderRepo := repositories.NewOrderRepository(db)
	productRepo := repositories.NewProductRepository(db)
	customerRepo := repositories.NewCustomerRepository(db)
	service := services.NewOrderService(orderRepo, productRepo, customerRepo, nil)

	t.Run("Create Order Success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO orders").WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec("INSERT INTO order_items").WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		order := &models.Order{
			CustomerID: "cust-1",
			Items: []models.OrderItem{
				{ProductID: "prod-1", Quantity: 1, Price: 10.99},
			},
		}

		err := service.CreateOrder(context.Background(), order)
		assert.NoError(t, err)
		assert.NotEmpty(t, order.ID)
		assert.Equal(t, 10.99, order.TotalAmount)
	})

	t.Run("Get Order Success", func(t *testing.T) {
		orderRows := sqlmock.NewRows([]string{"id", "customer_id", "total_amount"}).
			AddRow("order-1", "cust-1", 10.99)
		mock.ExpectQuery("SELECT id, customer_id, total_amount FROM orders").WillReturnRows(orderRows)

		itemRows := sqlmock.NewRows([]string{"id", "product_id", "quantity", "price"}).
			AddRow("item-1", "prod-1", 1, 10.99)
		mock.ExpectQuery("SELECT id, product_id, quantity, price FROM order_items").WillReturnRows(itemRows)

		order, err := service.GetOrder(context.Background(), "order-1")
		assert.NoError(t, err)
		assert.Equal(t, "order-1", order.ID)
		assert.Len(t, order.Items, 1)
	})
}
