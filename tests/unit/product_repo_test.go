package unit

import (
	"context"
	"testing"

	"ecommerce-service/internal/models"
	"ecommerce-service/internal/repositories"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestProductRepository(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock: %v", err)
	}
	defer db.Close()

	repo := repositories.NewProductRepository(db)

	t.Run("Create Product", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO products").
			WithArgs(sqlmock.AnyArg(), "Test Product", nil, 10.99, "cat-1").
			WillReturnResult(sqlmock.NewResult(1, 1))

		product := &models.Product{
			Name:       "Test Product",
			Price:      10.99,
			CategoryID: "cat-1",
		}

		err := repo.Create(context.Background(), product)
		assert.NoError(t, err)
		assert.NotEmpty(t, product.ID)
	})

	t.Run("Get Product By ID", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name", "price", "category_id"}).
			AddRow("prod-1", "Test Product", 10.99, "cat-1")
		mock.ExpectQuery("SELECT id, name, price, category_id FROM products WHERE id = ?").
			WithArgs("prod-1").
			WillReturnRows(rows)

		product, err := repo.GetByID(context.Background(), "prod-1")
		assert.NoError(t, err)
		assert.Equal(t, "prod-1", product.ID)
		assert.Equal(t, "Test Product", product.Name)
	})
}
