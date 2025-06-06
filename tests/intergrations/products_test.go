package integrations

import (
	"testing"

	"ecommerce-service/tests"

	"github.com/stretchr/testify/assert"
)

func TestProductIntegration(t *testing.T) {
	db := tests.SetupTestDB(t)

	t.Run("Database Operations", func(t *testing.T) {
		// Test direct database operations
		_, err := db.Exec("INSERT INTO products (id, name, price) VALUES (?, ?, ?)",
			"int-test-1", "Integration Product", 25.99)
		assert.NoError(t, err)

		var price float64
		err = db.QueryRow("SELECT price FROM products WHERE id = ?", "int-test-1").Scan(&price)
		assert.NoError(t, err)
		assert.Equal(t, 25.99, price)
	})
}
