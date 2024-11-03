package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/zde37/instashop-task/internal/models"
	"github.com/zde37/instashop-task/internal/repository"
	"github.com/zde37/instashop-task/pkg"
)

type serviceImpl struct {
	repo     repository.Repository
	jwtMaker *pkg.JWTMaker
}

func New(repo repository.Repository, jwtMaker *pkg.JWTMaker) Service {
	return &serviceImpl{
		repo:     repo,
		jwtMaker: jwtMaker,
	}
}

func (s *serviceImpl) Register(ctx context.Context, req *models.AuthRequest) (*models.User, error) {
	existingUser, err := s.repo.GetUser(ctx, pkg.EmailIdentifier, req.Email)
	if err != nil && !errors.Is(err, pkg.ErrNotFound) {
		return nil, fmt.Errorf("checking existing user: %w", err)
	}
	if existingUser != nil {
		return nil, pkg.ErrEmailTaken
	}

	hashedPassword, err := pkg.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("hashing password: %w", err)
	}

	user := &models.User{
		ID:           pkg.GenerateID(),
		Email:        req.Email,
		PasswordHash: hashedPassword,
		Role:         models.RoleCustomer,
	}

	if err := s.repo.CreateUser(ctx, user); err != nil {
		return nil, fmt.Errorf("creating user: %w", err)
	}
	return user, nil
}

func (s *serviceImpl) Login(ctx context.Context, req *models.AuthRequest) (accessToken, refreshToken string, err error) {
	user, err := s.repo.GetUser(ctx, pkg.EmailIdentifier, req.Email)
	if err != nil {
		if errors.Is(err, pkg.ErrNotFound) {
			return "", "", pkg.ErrInvalidCredentials
		}
		return "", "", fmt.Errorf("getting user: %w", err)
	}

	if err := pkg.VerifyPassword(user.PasswordHash, req.Password); err != nil {
		return "", "", pkg.ErrInvalidCredentials
	}

	// create tokens
	accessToken, err = s.jwtMaker.CreateToken(user.ID, user.Email, user.Role, pkg.AccessToken, pkg.AccessTokenDuration)
	if err != nil {
		return "", "", fmt.Errorf("creating access token: %w", err)
	}

	refreshToken, err = s.jwtMaker.CreateToken(user.ID, user.Email, user.Role, pkg.RefreshToken, pkg.RefreshTokenDuration)
	if err != nil {
		return "", "", fmt.Errorf("creating refresh token: %w", err)
	}

	// store refresh token in db
	session := &models.Session{
		ID:           pkg.GenerateID(),
		UserID:       user.ID,
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().Add(pkg.RefreshTokenDuration),
	}

	if err := s.repo.CreateSession(ctx, session); err != nil {
		return "", "", fmt.Errorf("storing session: %w", err)
	}

	return accessToken, refreshToken, nil
}

func (s *serviceImpl) Logout(ctx context.Context, userID string) error {
	return s.repo.DeleteSessionByUserID(ctx, userID)
}

func (s *serviceImpl) CreateProduct(ctx context.Context, req *models.CreateProductRequest) (*models.Product, error) {
	product := &models.Product{
		ID:            pkg.GenerateID(),
		Name:          req.Name,
		Description:   req.Description,
		Price:         req.Price,
		StockQuantity: req.StockQuantity,
	}

	if err := s.repo.CreateProduct(ctx, product); err != nil {
		return nil, fmt.Errorf("creating product: %w", err)
	}
	return product, nil
}

func (s *serviceImpl) GetProductByID(ctx context.Context, id string) (*models.Product, error) {
	product, err := s.repo.GetProductByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("getting product: %w", err)
	}
	return product, nil
}

func (s *serviceImpl) ListProducts(ctx context.Context) ([]models.Product, error) {
	products, err := s.repo.ListProducts(ctx)
	if err != nil {
		return nil, fmt.Errorf("listing products: %w", err)
	}

	return products, nil
}

func (s *serviceImpl) UpdateProduct(ctx context.Context, id string, req *models.CreateProductRequest) (*models.Product, error) {
	product := &models.Product{
		ID:            id,
		Name:          req.Name,
		Description:   req.Description,
		Price:         req.Price,
		StockQuantity: req.StockQuantity,
	}

	if err := s.repo.UpdateProduct(ctx, product); err != nil {
		return nil, fmt.Errorf("updating product: %w", err)
	}

	return product, nil
}

func (s *serviceImpl) DeleteProduct(ctx context.Context, id string) error {
	if err := s.repo.DeleteProduct(ctx, id); err != nil {
		return fmt.Errorf("deleting product: %w", err)
	}
	return nil
}

func (s *serviceImpl) CreateOrder(ctx context.Context, userID string, req *models.CreateOrderRequest) error {
	err := s.repo.WithTransaction(ctx, func(tx pgx.Tx) error {
		order := &models.Order{
			UserID: userID,
			Status: models.StatusPending,
			Items:  make([]models.OrderItem, len(req.Items)),
		}

		// process each order item
		for i, item := range req.Items {
			// get product
			product, err := s.repo.GetProductByID(ctx, item.ProductID)
			if err != nil {
				return fmt.Errorf("getting product %s: %w", item.ProductID, err)
			}

			// check stock
			if product.StockQuantity < item.Quantity {
				return pkg.ErrInsufficientStock
			}

			// update stock
			product.StockQuantity -= item.Quantity
			if err := s.repo.UpdateProduct(ctx, product); err != nil {
				return fmt.Errorf("updating product stock: %w", err)
			}

			// create order item
			order.Items[i] = models.OrderItem{
				ProductID: product.ID,
				Quantity:  item.Quantity,
				UnitPrice: product.Price,
				Product:   product,
			}

			// add to total
			order.TotalAmount += product.Price * float64(item.Quantity)
		}

		// create order
		err := s.repo.CreateOrder(ctx, tx, order)
		if err != nil {
			return fmt.Errorf("creating order: %w", err)
		}

		return err
	})
	return err
}

func (s *serviceImpl) GetOrderByID(ctx context.Context, id string) (*models.Order, error) {
	order, err := s.repo.GetOrderByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("getting order: %w", err)
	}
	return order, nil
}

func (s *serviceImpl) GetUserOrders(ctx context.Context, userID string) ([]models.Order, error) {
	orders, err := s.repo.GetOrderByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("getting user orders: %w", err)
	}
	return orders, nil
}

func (s *serviceImpl) UpdateStatus(ctx context.Context, id string, status models.OrderStatus) error {
	// get current order
	order, err := s.repo.GetOrderByID(ctx, id)
	if err != nil {
		return fmt.Errorf("getting order: %w", err)
	}

	// validate status transition
	if order.Status == models.StatusCancelled || order.Status == models.StatusDelivered {
		return fmt.Errorf("cannot update %s order", order.Status)
	}

	if err := s.repo.UpdateOrderStatus(ctx, id, status); err != nil {
		return fmt.Errorf("updating order status: %w", err)
	}
	return nil
}

func (s *serviceImpl) CancelOrder(ctx context.Context, id string, userID string) error {
	// get order
	order, err := s.repo.GetOrderByID(ctx, id)
	if err != nil {
		return fmt.Errorf("getting order: %w", err)
	}

	// verify ownership
	if order.UserID != userID {
		return pkg.ErrUnauthorized
	}

	// verify status
	if order.Status != models.StatusPending {
		return pkg.ErrOrderNotPending
	}

	// start transaction
	err = s.repo.WithTransaction(ctx, func(tx pgx.Tx) error {
		// update order status
		if err := s.repo.UpdateOrderStatus(ctx, id, models.StatusCancelled); err != nil {
			return fmt.Errorf("updating order status: %w", err)
		}

		// restore product stock
		for _, item := range order.Items {
			product, err := s.repo.GetProductByID(ctx, item.ProductID)
			if err != nil {
				return fmt.Errorf("getting product: %w", err)
			}
			product.StockQuantity += item.Quantity

			if err := s.repo.UpdateProduct(ctx, product); err != nil {
				return fmt.Errorf("restoring product stock: %w", err)
			}

		}
		return err
	})
	return err
}
