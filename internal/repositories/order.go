// internal/repositories/order.go
package repositories

import (
	"context"
	"database/sql"
	"ecommerce-service/internal/models"
	"time"
)

type OrderRepository interface {
	Create(ctx context.Context, order *models.Order) error
	GetByID(ctx context.Context, id string) (*models.Order, error)
	//PingDB(ctx context.Context) error
}

type orderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) Create(ctx context.Context, order *models.Order) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Insert order
	orderQuery := `INSERT INTO orders (id, customer_id, total_amount) VALUES (?, ?, ?)`
	_, err = tx.ExecContext(ctx, orderQuery, order.ID, order.CustomerID, order.TotalAmount)
	if err != nil {
		return err
	}

	// Insert order items
	itemQuery := `INSERT INTO order_items (id, order_id, product_id, quantity, price) VALUES (?, ?, ?, ?, ?)`
	for _, item := range order.Items {
		_, err = tx.ExecContext(ctx, itemQuery, item.ID, order.ID, item.ProductID, item.Quantity, item.Price)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *orderRepository) GetByID(ctx context.Context, id string) (*models.Order, error) {
	// Get order details
	orderQuery := `SELECT id, customer_id, order_date, status, total_amount 
	               FROM orders WHERE id = ?`
	row := r.db.QueryRowContext(ctx, orderQuery, id)

	var order models.Order
	var orderDate time.Time
	err := row.Scan(&order.ID, &order.CustomerID, &orderDate, &order.Status, &order.TotalAmount)
	if err != nil {
		return nil, err
	}
	order.OrderDate = orderDate

	// Get order items
	itemsQuery := `SELECT id, product_id, quantity, price 
	               FROM order_items WHERE order_id = ?`
	rows, err := r.db.QueryContext(ctx, itemsQuery, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.OrderItem
	for rows.Next() {
		var item models.OrderItem
		if err := rows.Scan(&item.ID, &item.ProductID, &item.Quantity, &item.Price); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	order.Items = items

	return &order, nil
}

// internal/repositories/order.go
func (r *orderRepository) Ping(ctx context.Context) error {
	return r.db.PingContext(ctx)
}
