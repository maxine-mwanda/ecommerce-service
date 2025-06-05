// internal/services/category.go
package services

import (
	"context"
	"ecommerce-service/internal/models"
	"ecommerce-service/internal/repositories"
)

type categoryService struct {
	repo repositories.CategoryRepository
}

func NewCategoryService(repo repositories.CategoryRepository) CategoryService {
	return &categoryService{repo: repo}
}

func (s *categoryService) CreateCategory(ctx context.Context, category *models.Category) error {
	// Validate parent exists if specified
	if category.ParentID != nil && *category.ParentID != "" {
		if _, err := s.repo.GetByID(ctx, *category.ParentID); err != nil {
			return err
		}
	}
	return s.repo.Create(ctx, category)
}

func (s *categoryService) GetCategory(ctx context.Context, id string) (*models.Category, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *categoryService) GetAveragePrice(ctx context.Context, categoryID string) (float64, error) {
	return s.repo.GetAveragePrice(ctx, categoryID)
}
