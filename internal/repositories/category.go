// internal/repositories/category.go
package repositories

import (
	"context"
	"database/sql"
	"ecommerce-service/internal/models"
)

type CategoryRepository struct {
	db *sql.DB
}

func NewCategoryRepository(db *sql.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) Create(ctx context.Context, category *models.Category) error {
	query := `INSERT INTO categories (id, name, parent_id) VALUES (?, ?, ?)`
	_, err := r.db.ExecContext(ctx, query, category.ID, category.Name, category.ParentID)
	return err
}

func (r *CategoryRepository) GetAveragePrice(ctx context.Context, categoryID string) (float64, error) {
	query := `SELECT AVG(price) FROM products WHERE category_id = ?`
	row := r.db.QueryRowContext(ctx, query, categoryID)

	var avgPrice sql.NullFloat64
	if err := row.Scan(&avgPrice); err != nil {
		return 0, err
	}

	if !avgPrice.Valid {
		return 0, nil
	}

	return avgPrice.Float64, nil
}
