package integration_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"ecommerce-service/internal/handlers"
	"ecommerce-service/internal/models"
	"ecommerce-service/internal/repositories"
	"ecommerce-service/internal/services"
	"ecommerce-service/tests"

	"github.com/stretchr/testify/assert"
)

func TestOrderIntegration(t *testing.T) {
	db := tests.SetupTestDB(t)

	// Setup test data
	_, err := db.Exec(`
		INSERT INTO customers (id, name, email) VALUES ('cust-1', 'Test User', 'test@example.com');
		INSERT INTO categories (id, name) VALUES ('cat-1', 'Test Category');
		INSERT INTO products (id, name, price, category_id) VALUES 
			('prod-1', 'Product 1', 10.99, 'cat-1'),
			('prod-2', 'Product 2', 20.50, 'cat-1');
	`)
	assert.NoError(t, err)

	// Initialize services
	orderRepo := repositories.NewOrderRepository(db)
	productRepo := repositories.NewProductRepository(db)
	customerRepo := repositories.NewCustomerRepository(db)
	orderService := services.NewOrderService(orderRepo, productRepo, customerRepo, nil)

	// Setup handler
	h := handlers.NewAPIHandler(
		nil, // auth service
		orderService,
		nil, // product service
		nil, // category service
		nil, // logger
	)

	t.Run("Create Order", func(t *testing.T) {
		order := models.Order{
			CustomerID: "cust-1",
			Items: []models.OrderItem{
				{ProductID: "prod-1", Quantity: 2, Price: 10.99},
				{ProductID: "prod-2", Quantity: 1, Price: 20.50},
			},
		}

		body, _ := json.Marshal(order)
		req := httptest.NewRequest("POST", "/orders", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response models.Order
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.NotEmpty(t, response.ID)
		assert.Equal(t, 42.48, response.TotalAmount) // 10.99*2 + 20.50
		assert.Len(t, response.Items, 2)
	})
}
