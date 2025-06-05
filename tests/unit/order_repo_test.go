package repositories_test

import (
	"context"
	"testing"
	"time"

	"ecommerce-service/internal/models"
	"ecommerce-service/internal/repositories"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestOrderRepository(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock: %v", err)
	}
	defer db.Close()

	repo := repositories.NewOrderRepository(db)

	t.Run("Create Order", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO orders").
			WithArgs("order-1", "cust-1", 10.99).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec("INSERT INTO order_items").
			WithArgs("item-1", "order-1", "prod-1", 1, 10.99).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		order := &models.Order{
			ID:          "order-1",
			CustomerID:  "cust-1",
			TotalAmount: 10.99,
			Items: []models.OrderItem{
				{ID: "item-1", ProductID: "prod-1", Quantity: 1, Price: 10.99},
			},
		}

		err := repo.Create(context.Background(), order)
		assert.NoError(t, err)
	})

	t.Run("Get Order By ID", func(t *testing.T) {
		orderTime := time.Now()
		orderRows := sqlmock.NewRows([]string{"id", "customer_id", "order_date", "total_amount"}).
			AddRow("order-1", "cust-1", orderTime, 10.99)
		mock.ExpectQuery("SELECT id, customer_id, order_date, total_amount FROM orders WHERE id = ?").
			WithArgs("order-1").
			WillReturnRows(orderRows)

		itemRows := sqlmock.NewRows([]string{"id", "product_id", "quantity", "price"}).
			AddRow("item-1", "prod-1", 1, 10.99)
		mock.ExpectQuery("SELECT id, product_id, quantity, price FROM order_items WHERE order_id = ?").
			WithArgs("order-1").
			WillReturnRows(itemRows)

		order, err := repo.GetByID(context.Background(), "order-1")
		assert.NoError(t, err)
		assert.Equal(t, "order-1", order.ID)
		assert.Len(t, order.Items, 1)
		assert.Equal(t, "prod-1", order.Items[0].ProductID)
	})
}
