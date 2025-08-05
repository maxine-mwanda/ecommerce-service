package unit

import (
	"context"
	"testing"

	"ecommerce-service/internal/models"
	"ecommerce-service/internal/repositories"
	"ecommerce-service/tests"

	"github.com/stretchr/testify/assert"
)

func TestProductRepository(t *testing.T) {
	db := tests.SetupTestDB(t)
	repo := repositories.NewProductRepository(db)

	t.Run("Create and Get Product", func(t *testing.T) {
		// Create test product
		product := &models.Product{
			ID:         "test-product",
			Name:       "Test Product",
			Price:      10.99,
			CategoryID: "cat-1",
		}

		err := repo.Create(context.Background(), product)
		assert.NoError(t, err)

		// Retrieve product
		retrieved, err := repo.GetByID(context.Background(), "test-product")
		assert.NoError(t, err)
		assert.Equal(t, "Test Product", retrieved.Name)
	})
}
