// internal/services/product.go
package services

import (
	"context"
	"ecommerce-service/internal/models"
	"ecommerce-service/internal/repositories"

	"github.com/google/uuid"
)

type productService struct {
	productRepo  repositories.ProductRepository
	categoryRepo repositories.CategoryRepository
}

func NewProductService(
	productRepo repositories.ProductRepository,
	categoryRepo repositories.CategoryRepository,
) ProductService {
	return &productService{
		productRepo:  productRepo,
		categoryRepo: categoryRepo,
	}
}

func (s *productService) CreateProduct(ctx context.Context, product *models.Product) error {
	// Verify category exists
	if product.CategoryID == "" {
		product.ID = uuid.New().String()
	}
	if _, err := s.categoryRepo.GetByID(ctx, product.CategoryID); err != nil {
		return err
	}
	return s.productRepo.Create(ctx, product)
}

func (s *productService) GetProduct(ctx context.Context, id string) (*models.Product, error) {
	return s.productRepo.GetByID(ctx, id)
}

func (s *productService) GetProductsByCategory(ctx context.Context, categoryID string) ([]models.Product, error) {
	return s.productRepo.GetByCategory(ctx, categoryID)
}
