package models

type AuthRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,validpassword"`
}

type CreateProductRequest struct {
	Name          string  `json:"name" binding:"required"`
	Description   string  `json:"description"`
	Price         float64 `json:"price" binding:"required,gt=0"`
	StockQuantity int     `json:"stock_quantity" binding:"required,gte=0"`
}

type CreateOrderRequest struct {
	Items []CreateOrderItemRequest `json:"items" binding:"required,dive"`
}

type CreateOrderItemRequest struct {
	ProductID string `json:"product_id" binding:"required"`
	Quantity  int    `json:"quantity" binding:"required,gt=0"`
}

type UpdateOrderStatusRequest struct {
	Status OrderStatus `json:"status" binding:"required,oneof=pending confirmed shipped delivered cancelled"`
}

// ErrorResponse represents the error response structure
type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
}
