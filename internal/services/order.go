// internal/services/order.go
package services

import (
	"context"
	"ecommerce-service/internal/models"
	"ecommerce-service/internal/repositories"
	"fmt"
	"strconv"
)

type orderService struct {
	orderRepo    repositories.OrderRepository
	productRepo  repositories.ProductRepository
	customerRepo repositories.CustomerRepository
	notifier     NotificationService
}

func NewOrderService(
	orderRepo repositories.OrderRepository,
	productRepo repositories.ProductRepository,
	customerRepo repositories.CustomerRepository,
	notifier NotificationService,
) OrderService {
	return &orderService{
		orderRepo:    orderRepo,
		productRepo:  productRepo,
		customerRepo: customerRepo,
		notifier:     notifier,
	}
}

func (s *orderService) CreateOrder(ctx context.Context, order *models.Order) error {
	// Calculate total amount
	var total float64
	for _, item := range order.Items {
		total += item.Price * float64(item.Quantity)
	}
	order.TotalAmount = total

	// Create order in database
	if err := s.orderRepo.Create(ctx, order); err != nil {
		return err
	}
	if err := s.orderRepo.Create(ctx, order); err != nil {
		return fmt.Errorf("orderService.CreateOrder failed: %w", err)
	}

	// Get customer details for notifications
	customer, err := s.customerRepo.GetByEmail(ctx, order.CustomerID)
	if err != nil {
		return err
	}

	// Send notifications (in background)
	go s.notifier.SendOrderConfirmationSMS(customer.Phone, "Your order #"+order.ID+" has been placed")
	go s.notifier.SendOrderNotificationEmail(
		"admin@example.com",
		"New Order #"+order.ID,
		"Customer: "+customer.Name+"\nTotal: $"+strconv.FormatFloat(order.TotalAmount, 'f', 2, 64),
	)

	return nil
}

func (s *orderService) GetOrder(ctx context.Context, id string) (*models.Order, error) {
	return s.orderRepo.GetByID(ctx, id)
}
