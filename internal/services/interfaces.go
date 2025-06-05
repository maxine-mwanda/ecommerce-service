// internal/services/interfaces.go
package services

import (
	"context"
	"ecommerce-service/internal/models"
)

type AuthService interface {
	Authenticate(token string) (string, error)
}

type OrderService interface {
	CreateOrder(ctx context.Context, order *models.Order) error
	GetOrder(ctx context.Context, id string) (*models.Order, error)
}

type ProductService interface {
	CreateProduct(ctx context.Context, product *models.Product) error
	GetProduct(ctx context.Context, id string) (*models.Product, error)
	GetProductsByCategory(ctx context.Context, categoryID string) ([]models.Product, error)
}

type CategoryService interface {
	CreateCategory(ctx context.Context, category *models.Category) error
	GetCategory(ctx context.Context, id string) (*models.Category, error)
	GetAveragePrice(ctx context.Context, categoryID string) (float64, error)
}

/*type NotificationService interface {
	SendOrderConfirmationSMS(phone, message string) error
	SendOrderNotificationEmail(email, subject, message string) error
}*/
