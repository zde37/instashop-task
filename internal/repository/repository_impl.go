package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/zde37/instashop-task/internal/models"
	"github.com/zde37/instashop-task/pkg"
)

type repositoryImpl struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) Repository {
	return &repositoryImpl{db: db}
}

func (r *repositoryImpl) CreateUser(ctx context.Context, user *models.User) error {
	query := `INSERT INTO users (id, email, password_hash, role) VALUES ($1, $2, $3, $4) RETURNING created_at, updated_at`

	err := pgxscan.Get(ctx, r.db, user, query, user.ID, user.Email, user.PasswordHash, user.Role)
	if err != nil {
		return fmt.Errorf("create user: %w", err)
	}
	return nil
}

func (r *repositoryImpl) GetUser(ctx context.Context, identifier, data string) (*models.User, error) {
	var user models.User
	query := fmt.Sprintf(`SELECT id, email, password_hash, role, created_at, updated_at FROM users WHERE %s = $1`, identifier)

	err := pgxscan.Get(ctx, r.db, &user, query, data)
	if err != nil {
		if pgxscan.NotFound(err) {
			return nil, pkg.ErrNotFound
		}
		return nil, fmt.Errorf("get user by %s: %w", identifier, err)
	}
	return &user, nil
}

func (r *repositoryImpl) CreateProduct(ctx context.Context, product *models.Product) error {
	query := `INSERT INTO products (id, name, description, price, stock_quantity) VALUES ($1, $2, $3, $4, $5) RETURNING created_at, updated_at`

	err := pgxscan.Get(ctx, r.db, product, query, product.ID, product.Name, product.Description, product.Price, product.StockQuantity)
	if err != nil {
		return fmt.Errorf("create product: %w", err)
	}
	return nil
}

func (r *repositoryImpl) GetProductByID(ctx context.Context, id string) (*models.Product, error) {
	var product models.Product
	query := `SELECT id, name, description, price, stock_quantity, created_at, updated_at FROM products WHERE id = $1`

	err := pgxscan.Get(ctx, r.db, &product, query, id)
	if err != nil {
		if pgxscan.NotFound(err) {
			return nil, pkg.ErrNotFound
		}
		return nil, fmt.Errorf("get product by id: %w", err)
	}
	return &product, nil
}

func (r *repositoryImpl) ListProducts(ctx context.Context) ([]models.Product, error) {
	var products []models.Product
	query := `SELECT id, name, description, price, stock_quantity, created_at, updated_at FROM products ORDER BY created_at DESC`

	err := pgxscan.Select(ctx, r.db, &products, query)
	if err != nil {
		return nil, fmt.Errorf("list products: %w", err)
	}
	return products, nil
}

func (r *repositoryImpl) UpdateProduct(ctx context.Context, product *models.Product) error {
	query := `UPDATE products SET name = $1, description = $2, price = $3, stock_quantity = $4, updated_at = $5 WHERE id = $6 RETURNING updated_at`
	now := time.Now()

	err := pgxscan.Get(ctx, r.db, &product.UpdatedAt, query, product.Name, product.Description, product.Price, product.StockQuantity, now, product.ID)
	if err != nil {
		if pgxscan.NotFound(err) {
			return pkg.ErrNotFound
		}
		return fmt.Errorf("update product: %w", err)
	}
	return nil
}

func (r *repositoryImpl) DeleteProduct(ctx context.Context, id string) error {
	query := `DELETE FROM products WHERE id = $1`

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete product: %w", err)
	}

	if result.RowsAffected() == 0 {
		return pkg.ErrNotFound
	}
	return nil
}

func (r *repositoryImpl) CreateOrder(ctx context.Context, tx pgx.Tx, order *models.Order) error {
	query := `INSERT INTO orders (id, user_id, status, total_amount) VALUES ($1, $2, $3, $4) RETURNING created_at, updated_at`

	err := pgxscan.Get(ctx, tx, order, query, order.ID, order.UserID, order.Status, order.TotalAmount)
	if err != nil {
		return fmt.Errorf("create order: %w", err)
	}

	// insert order items
	for i := range order.Items {
		item := &order.Items[i]
		item.ID = pkg.GenerateID()
		item.OrderID = order.ID

		query = `INSERT INTO order_items (id, order_id, product_id, quantity, unit_price) VALUES ($1, $2, $3, $4, $5) RETURNING created_at, updated_at, subtotal`

		err = pgxscan.Get(ctx, tx, item, query, item.ID, item.OrderID, item.ProductID, item.Quantity, item.UnitPrice)
		if err != nil {
			return fmt.Errorf("create order item: %w", err)
		}
	}
	return err
}

func (r *repositoryImpl) WithTransaction(ctx context.Context, fn func(pgx.Tx) error) error {
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}

	err = fn(tx)
	if err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("tx err: %w, rb err: %w", err, rbErr)
		}
		return err
	}

	return tx.Commit(ctx)
}

func (r *repositoryImpl) GetOrderByID(ctx context.Context, id string) (*models.Order, error) {
	var order models.Order
	query := `SELECT id, user_id, status, total_amount, created_at, updated_at FROM orders WHERE id = $1`

	err := pgxscan.Get(ctx, r.db, &order, query, id)
	if err != nil {
		if pgxscan.NotFound(err) {
			return nil, pkg.ErrNotFound
		}
		return nil, fmt.Errorf("get order by id: %w", err)
	}

	// get order items with products
	query = `
        SELECT i.*, p.*
        FROM order_items i
        JOIN products p ON p.id = i.product_id
        WHERE i.order_id = $1`

	var items []struct {
		models.OrderItem
		Product models.Product `db:"product"`
	}

	err = pgxscan.Select(ctx, r.db, &items, query, id)
	if err != nil {
		return nil, fmt.Errorf("get order items: %w", err)
	}

	order.Items = make([]models.OrderItem, len(items))
	for i, item := range items {
		order.Items[i] = item.OrderItem
		order.Items[i].Product = &item.Product
	}

	return &order, nil
}

func (r *repositoryImpl) GetOrderByUserID(ctx context.Context, userID string) ([]models.Order, error) {
	query := `SELECT id, user_id, status, total_amount, created_at, updated_at FROM orders WHERE user_id = $1 ORDER BY created_at DESC`

	var orders []models.Order
	err := pgxscan.Select(ctx, r.db, &orders, query, userID)
	if err != nil {
		return nil, fmt.Errorf("get orders by user id: %w", err)
	}
	return orders, nil
}

func (r *repositoryImpl) UpdateOrderStatus(ctx context.Context, id string, status models.OrderStatus) error {
	query := `UPDATE orders SET status = $1, updated_at = $2 WHERE id = $3 RETURNING updated_at`
	now := time.Now()

	var updatedAt time.Time
	err := pgxscan.Get(ctx, r.db, &updatedAt, query, status, now, id)
	if err != nil {
		if pgxscan.NotFound(err) {
			return pkg.ErrNotFound
		}
		return fmt.Errorf("update order status: %w", err)
	}
	return nil
}

func (r *repositoryImpl) CreateSession(ctx context.Context, session *models.Session) error {
	query := `INSERT INTO sessions (id, user_id, refresh_token, expires_at) VALUES ($1, $2, $3, $4) RETURNING created_at`

	err := pgxscan.Get(ctx, r.db, session, query, session.ID, session.UserID, session.RefreshToken, session.ExpiresAt)
	if err != nil {
		return fmt.Errorf("create session: %w", err)
	}
	return nil
}

func (r *repositoryImpl) GetSessionByID(ctx context.Context, id string) (*models.Session, error) {
	var session models.Session
	query := `SELECT id, user_id, refresh_token, is_blocked, expires_at, created_at FROM sessions WHERE id = $1`

	err := pgxscan.Get(ctx, r.db, &session, query, id)
	if err != nil {
		if pgxscan.NotFound(err) {
			return nil, pkg.ErrNotFound
		}
		return nil, fmt.Errorf("get session by id: %w", err)
	}

	return &session, nil
}

func (r *repositoryImpl) DeleteSessionByUserID(ctx context.Context, userID string) error {
	query := `DELETE FROM sessions WHERE user_id = $1`

	result, err := r.db.Exec(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("delete sessions by user id: %w", err)
	}

	if result.RowsAffected() == 0 {
		return pkg.ErrNotFound
	}
	return nil
}

func (r *repositoryImpl) BlockSession(ctx context.Context, id string) error {
	query := `UPDATE sessions SET is_blocked = true WHERE id = $1`

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("block session: %w", err)
	}

	if result.RowsAffected() == 0 {
		return pkg.ErrNotFound
	}
	return nil
}
