// internal/services/order.go
package services

import (
	"context"
	"database/sql"
	"ecommerce-service/internal/models"
	"ecommerce-service/internal/repositories"
	"errors"
	"fmt"
	"strconv"

	"github.com/google/uuid"
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
	//avoiding empty ID
	if order.ID == "" {
		order.ID = generateUUID()
	}
	// Calculate total amount
	var total float64
	for _, item := range order.Items {
		total += item.Price * float64(item.Quantity)
	}
	order.TotalAmount = total

	// Check if customer exists by email
	customer, err := s.customerRepo.GetByEmail(ctx, order.CustomerEmail)
	if err != nil {
		// If not found, create a new customer
		if errors.Is(err, sql.ErrNoRows) {
			newCustomer := &models.Customer{
				ID:    generateUUID(),     // or whatever your ID generator is
				Name:  order.CustomerName, // get from order if available
				Email: order.CustomerEmail,
				Phone: order.CustomerPhone,
			}
			err = s.customerRepo.Create(ctx, newCustomer)
			if err != nil {
				return fmt.Errorf("failed to create customer: %w", err)
			}
			customer = newCustomer
		} else {
			return fmt.Errorf("failed to get customer: %w", err)
		}
	}

	// Set the valid customer ID on order
	order.CustomerID = customer.ID

	// Create order in database
	if err := s.orderRepo.Create(ctx, order); err != nil {
		return fmt.Errorf("orderService.CreateOrder failed: %w", err)
	}

	// Send notifications (in background)
	go s.notifier.SendOrderConfirmationSMS(customer.Phone, "Your order #"+order.ID+" has been placed")
	go s.notifier.SendOrderNotificationEmail(
		"mwandamaxyine@gmail.com",
		"New Order #"+order.ID,
		"Customer: "+customer.Name+"\nTotal: $"+strconv.FormatFloat(order.TotalAmount, 'f', 2, 64),
	)

	return nil
}

func generateUUID() string {
	return uuid.New().String()
}

func (s *orderService) GetOrder(ctx context.Context, id string) (*models.Order, error) {
	return s.orderRepo.GetByID(ctx, id)
}
