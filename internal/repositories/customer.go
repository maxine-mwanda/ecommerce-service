// internal/repositories/customer.go
package repositories

import (
	"context"
	"database/sql"
	"ecommerce-service/internal/models"
)

type CustomerRepository struct {
	db *sql.DB
}

func NewCustomerRepository(db *sql.DB) *CustomerRepository {
	return &CustomerRepository{db: db}
}

func (r *CustomerRepository) Create(ctx context.Context, customer *models.Customer) error {
	query := `INSERT INTO customers (id, name, email, phone) VALUES (?, ?, ?, ?)`
	_, err := r.db.ExecContext(ctx, query, customer.ID, customer.Name, customer.Email, customer.Phone)
	return err
}

func (r *CustomerRepository) GetByEmail(ctx context.Context, email string) (*models.Customer, error) {
	query := `SELECT id, name, email, phone, created_at, updated_at FROM customers WHERE email = ?`
	row := r.db.QueryRowContext(ctx, query, email)

	var customer models.Customer
	err := row.Scan(&customer.ID, &customer.Name, &customer.Email, &customer.Phone, &customer.CreatedAt, &customer.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &customer, nil
}
