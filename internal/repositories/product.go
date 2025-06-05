// internal/repositories/product.go
package repositories

import (
	"context"
	"database/sql"
	"ecommerce-service/internal/models"
)

type ProductRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) Create(ctx context.Context, product *models.Product) error {
	query := `INSERT INTO products (id, name, description, price, category_id) VALUES (?, ?, ?, ?, ?)`
	_, err := r.db.ExecContext(ctx, query,
		product.ID,
		product.Name,
		product.Description,
		product.Price,
		product.CategoryID,
	)
	return err
}

func (r *ProductRepository) GetByCategory(ctx context.Context, categoryID string) ([]models.Product, error) {
	query := `SELECT id, name, description, price, category_id, created_at, updated_at 
	          FROM products WHERE category_id = ?`
	rows, err := r.db.QueryContext(ctx, query, categoryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.CategoryID, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	return products, nil
}

// internal/repositories/product.go (add this method)
func (r *ProductRepository) GetByID(ctx context.Context, id string) (*models.Product, error) {
	query := `SELECT id, name, description, price, category_id, created_at, updated_at 
	          FROM products WHERE id = ?`
	row := r.db.QueryRowContext(ctx, query, id)

	var product models.Product
	err := row.Scan(
		&product.ID,
		&product.Name,
		&product.Description,
		&product.Price,
		&product.CategoryID,
		&product.CreatedAt,
		&product.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &product, nil
}
