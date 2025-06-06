package unit

import (
	"context"
	"testing"

	"ecommerce-service/internal/models"
	"ecommerce-service/internal/repositories"
	"ecommerce-service/internal/services"
	"ecommerce-service/tests"

	"github.com/stretchr/testify/assert"
)

func TestProductService(t *testing.T) {
	db := tests.SetupTestDB(t)
	repo := repositories.NewProductRepository(db)
	var productRepo repositories.ProductRepository = *repo
	service := services.NewProductService(productRepo, *repositories.NewCategoryRepository(db))

	t.Run("Create Product", func(t *testing.T) {
		product := &models.Product{
			Name:  "Service Test",
			Price: 15.99,
		}

		err := service.CreateProduct(context.Background(), product)
		assert.NoError(t, err)
		assert.NotEmpty(t, product.ID)
	})
}
