package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/zde37/instashop-task/internal/models"
)

type Repository interface {
	CreateUser(ctx context.Context, user *models.User) error
	GetUser(ctx context.Context, identifier, data string) (*models.User, error)

	WithTransaction(ctx context.Context, fn func(pgx.Tx) error) error
	CreateProduct(ctx context.Context, product *models.Product) error
	GetProductByID(ctx context.Context, id string) (*models.Product, error)
	ListProducts(ctx context.Context) ([]models.Product, error)
	UpdateProduct(ctx context.Context, product *models.Product) error
	DeleteProduct(ctx context.Context, id string) error

	CreateOrder(ctx context.Context, tx pgx.Tx, order *models.Order) error
	GetOrderByID(ctx context.Context, id string) (*models.Order, error)
	GetOrderByUserID(ctx context.Context, userID string) ([]models.Order, error)
	UpdateOrderStatus(ctx context.Context, id string, status models.OrderStatus) error

	CreateSession(ctx context.Context, session *models.Session) error
	GetSessionByID(ctx context.Context, id string) (*models.Session, error)
	DeleteSessionByUserID(ctx context.Context, userID string) error
	BlockSession(ctx context.Context, id string) error
}
