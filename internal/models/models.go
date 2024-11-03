package models

import (
	"time"
)

type UserRole string
type OrderStatus string

const (
	RoleCustomer UserRole = "customer"
	RoleAdmin    UserRole = "admin"

	StatusPending   OrderStatus = "pending"
	StatusConfirmed OrderStatus = "confirmed"
	StatusShipped   OrderStatus = "shipped"
	StatusDelivered OrderStatus = "delivered"
	StatusCancelled OrderStatus = "cancelled"
)

type User struct {
	ID           string    `json:"id" db:"id"`
	Email        string    `json:"email" db:"email"`
	PasswordHash string    `json:"-" db:"password_hash"`
	Role         UserRole  `json:"role" db:"role"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

type Product struct {
	ID            string    `json:"id" db:"id"`
	Name          string    `json:"name" db:"name"`
	Description   string    `json:"description" db:"description"`
	Price         float64   `json:"price" db:"price"`
	StockQuantity int       `json:"stock_quantity" db:"stock_quantity"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}

type Order struct {
	ID          string      `json:"id" db:"id"`
	UserID      string      `json:"user_id" db:"user_id"`
	Status      OrderStatus `json:"status" db:"status"`
	TotalAmount float64     `json:"total_amount" db:"total_amount"`
	CreatedAt   time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at" db:"updated_at"`
	Items       []OrderItem `json:"items,omitempty" db:"-"`
}

type OrderItem struct {
	ID        string    `json:"id" db:"id"`
	OrderID   string    `json:"order_id" db:"order_id"`
	ProductID string    `json:"product_id" db:"product_id"`
	Quantity  int       `json:"quantity" db:"quantity"`
	UnitPrice float64   `json:"unit_price" db:"unit_price"`
	SubTotal  float64   `json:"sub_total" db:"sub_total"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	Product   *Product  `json:"product,omitempty" db:"-"`
}

type Session struct {
	ID           string    `json:"id" db:"id"`
	UserID       string    `json:"user_id" db:"user_id"`
	RefreshToken string    `json:"refresh_token" db:"refresh_token"`
	IsBlocked    bool      `json:"is_blocked" db:"is_blocked"`
	ExpiresAt    time.Time `json:"expires_at" db:"expires_at"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}
