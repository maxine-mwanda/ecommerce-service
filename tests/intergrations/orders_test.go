package integrations

import (
	"testing"

	"ecommerce-service/tests"

	"github.com/stretchr/testify/assert"
)

func TestOrderIntegration(t *testing.T) {
	db := tests.SetupTestDB(t)

	// Setup test data
	_, err := db.Exec(`
		INSERT INTO products (id, name, price) VALUES 
		('prod-1', 'Product 1', 10.50),
		('prod-2', 'Product 2', 20.00);
	`)
	assert.NoError(t, err)

	t.Run("Order Creation", func(t *testing.T) {
		_, err := db.Exec(`
			INSERT INTO orders (id, customer_id, total) VALUES
			('order-1', 'cust-1', 30.50);
		`)
		assert.NoError(t, err)
	})
}
