// internal/models/models.go
package models

import (
	"time"
)

type Customer struct {
	ID        string    `json:"id"`
	Name      string    `json:"name" validate:"required"`
	Email     string    `json:"email" validate:"required,email"`
	Phone     string    `json:"phone"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Category struct {
	ID       string     `json:"id"`
	Name     string     `json:"name" validate:"required"`
	ParentID *string    `json:"parent_id"`
	Children []Category `json:"children,omitempty"`
}

type Product struct {
	ID          string    `json:"id"`
	Name        string    `json:"name" validate:"required"`
	Description string    `json:"description"`
	Price       float64   `json:"price" validate:"required,gt=0"`
	CategoryID  string    `json:"category_id" validate:"required"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Order struct {
	ID          string      `json:"id"`
	CustomerID  string      `json:"customer_id" validate:"required"`
	OrderDate   time.Time   `json:"order_date"`
	Status      string      `json:"status"`
	TotalAmount float64     `json:"total_amount"`
	Items       []OrderItem `json:"items"`
}

type OrderItem struct {
	ID        string  `json:"id"`
	OrderID   string  `json:"order_id"`
	ProductID string  `json:"product_id" validate:"required"`
	Quantity  int     `json:"quantity" validate:"required,gt=0"`
	Price     float64 `json:"price" validate:"required,gt=0"`
}
