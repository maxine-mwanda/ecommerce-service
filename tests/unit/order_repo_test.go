package unit

import (
	"context"
	"testing"

	"ecommerce-service/internal/models"
	"ecommerce-service/internal/repositories"
	"ecommerce-service/tests"

	"github.com/stretchr/testify/assert"
)

func TestOrderRepository(t *testing.T) {
	db := tests.SetupTestDB(t)
	repo := repositories.NewOrderRepository(db)

	t.Run("Create and Get Order", func(t *testing.T) {
		// Create test order
		order := &models.Order{
			ID:          "test-order",
			CustomerID:  "test-customer",
			TotalAmount: 20.99,
		}

		err := repo.Create(context.Background(), order)
		assert.NoError(t, err)

		// Retrieve order
		retrieved, err := repo.GetByID(context.Background(), "test-order")
		assert.NoError(t, err)
		assert.Equal(t, "test-customer", retrieved.CustomerID)
	})
}
