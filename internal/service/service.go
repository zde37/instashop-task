package service

import (
	"context"

	"github.com/zde37/instashop-task/internal/models"
)

type Service interface {
	Register(ctx context.Context, req *models.AuthRequest) (*models.User, error)
	Login(ctx context.Context, req *models.AuthRequest) (accessToken, refreshToken string, err error)
	Logout(ctx context.Context, userID string) error

	CreateProduct(ctx context.Context, req *models.CreateProductRequest) (*models.Product, error)
	GetProductByID(ctx context.Context, id string) (*models.Product, error)
	ListProducts(ctx context.Context) ([]models.Product, error)
	UpdateProduct(ctx context.Context, id string, req *models.CreateProductRequest) (*models.Product, error)
	DeleteProduct(ctx context.Context, id string) error

	CreateOrder(ctx context.Context, userID string, req *models.CreateOrderRequest) error
	GetOrderByID(ctx context.Context, id string) (*models.Order, error)
	GetUserOrders(ctx context.Context, userID string) ([]models.Order, error)
	UpdateStatus(ctx context.Context, id string, status models.OrderStatus) error
	CancelOrder(ctx context.Context, id string, userID string) error
}
