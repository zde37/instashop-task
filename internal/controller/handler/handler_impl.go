package handler

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zde37/instashop-task/internal/models"
	"github.com/zde37/instashop-task/internal/service"
	"github.com/zde37/instashop-task/pkg"
)

type handlerImpl struct {
	service service.Service
}

func New(service service.Service) Handler {
	return &handlerImpl{
		service: service,
	}
}

// Register
// @Summary      Register a new user
// @Description  Register a new user with email and password
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body models.AuthRequest true "User registration credentials"
// @Success      201 {object} models.User
// @Failure      400 {object} models.ErrorResponse
// @Failure      409 {object} models.ErrorResponse
// @Failure      500 {object} models.ErrorResponse
// @Router       /register [post]
func (h *handlerImpl) Register(ctx *gin.Context) {
	var req models.AuthRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.handleError(ctx, pkg.ErrInvalidInput, "register_validation")
		return
	}

	user, err := h.service.Register(ctx.Request.Context(), &req)
	if err != nil {
		h.handleError(ctx, err, "register_user")
		return
	}
	ctx.JSON(http.StatusCreated, user)
}

// Login
// @Summary      Login user
// @Description  Authenticate user and return access and refresh tokens
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body models.AuthRequest true "User login credentials"
// @Success      200 {object} map[string]string
// @Failure      400 {object} models.ErrorResponse
// @Failure      401 {object} models.ErrorResponse
// @Failure      500 {object} models.ErrorResponse
// @Router       /login [post]
func (h *handlerImpl) Login(c *gin.Context) {
	var req models.AuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.handleError(c, pkg.ErrInvalidInput, "login_validation")
		return
	}

	accessToken, refreshToken, err := h.service.Login(c.Request.Context(), &req)
	if err != nil {
		h.handleError(c, err, "login_user")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

// Logout
// @Summary      Logout user
// @Description  Invalidate the user's refresh token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Success      204 "No Content"
// @Failure      401 {object} models.ErrorResponse
// @Failure      500 {object} models.ErrorResponse
// @Security     Bearer
// @Router       /logout [post]
func (h *handlerImpl) Logout(c *gin.Context) {
	userID := c.GetString("user_id")
	if err := h.service.Logout(c.Request.Context(), userID); err != nil {
		h.handleError(c, err, "logout")
		return
	}

	c.Status(http.StatusNoContent)
}

// CreateProduct
// @Summary      Create a new product
// @Description  Create a new product (admin only)
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        request body models.CreateProductRequest true "Product details"
// @Success      201 {object} models.Product
// @Failure      400 {object} models.ErrorResponse
// @Failure      401 {object} models.ErrorResponse
// @Failure      403 {object} models.ErrorResponse
// @Failure      500 {object} models.ErrorResponse
// @Security     Bearer
// @Router       /products [post]
func (h *handlerImpl) CreateProduct(c *gin.Context) {
	var req models.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.handleError(c, pkg.ErrInvalidInput, "create_product_validation")
		return
	}

	product, err := h.service.CreateProduct(c.Request.Context(), &req)
	if err != nil {
		h.handleError(c, err, "create_product")
		return
	}

	c.JSON(http.StatusCreated, product)
}

// GetProduct
// @Summary      Get product by ID
// @Description  Get detailed information about a specific product
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        id path string true "Product ID"
// @Success      201 {object} models.Product
// @Failure      400 {object} models.ErrorResponse
// @Failure      401 {object} models.ErrorResponse
// @Failure      403 {object} models.ErrorResponse
// @Failure      500 {object} models.ErrorResponse
// @Security     Bearer
// @Router       /products/{id} [get]
func (h *handlerImpl) GetProduct(c *gin.Context) {
	id := c.Param("id")

	product, err := h.service.GetProductByID(c.Request.Context(), id)
	if err != nil {
		h.handleError(c, err, "get_product")
		return
	}

	c.JSON(http.StatusOK, product)
}

// ListProducts
// @Summary      List all products
// @Description  Get a list of all available products
// @Tags         products
// @Accept       json
// @Produce      json
// @Success      200 {array} models.Product
// @Failure      400 {object} models.ErrorResponse
// @Failure      401 {object} models.ErrorResponse
// @Failure      403 {object} models.ErrorResponse
// @Failure      500 {object} models.ErrorResponse
// @Security     Bearer
// @Router       /products [get]
func (h *handlerImpl) ListProducts(c *gin.Context) {
	products, err := h.service.ListProducts(c.Request.Context())
	if err != nil {
		h.handleError(c, err, "list_products")
		return
	}

	c.JSON(http.StatusOK, products)
}

// UpdateProduct
// @Summary      Update product
// @Description  Update product details (admin only)
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        id path string true "Product ID"
// @Param        request body models.CreateProductRequest true "Product details"
// @Success      200 {object} models.Product
// @Failure      400 {object} models.ErrorResponse
// @Failure      401 {object} models.ErrorResponse
// @Failure      403 {object} models.ErrorResponse
// @Failure      404 {object} models.ErrorResponse
// @Failure      500 {object} models.ErrorResponse
// @Security     Bearer
// @Router       /products/{id} [put]
func (h *handlerImpl) UpdateProduct(c *gin.Context) {
	id := c.Param("id")

	var req models.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.handleError(c, pkg.ErrInvalidInput, "update_product_validation")
		return
	}

	product, err := h.service.UpdateProduct(c.Request.Context(), id, &req)
	if err != nil {
		h.handleError(c, err, "update_product")
		return
	}

	c.JSON(http.StatusOK, product)
}

// DeleteProduct
// @Summary      Delete product
// @Description  Delete a product (admin only)
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        id path string true "Product ID"
// @Success      204 "No Content"
// @Failure      401 {object} models.ErrorResponse
// @Failure      403 {object} models.ErrorResponse
// @Failure      404 {object} models.ErrorResponse
// @Failure      500 {object} models.ErrorResponse
// @Security     Bearer
// @Router       /products/{id} [delete]
func (h *handlerImpl) DeleteProduct(c *gin.Context) {
	id := c.Param("id")

	if err := h.service.DeleteProduct(c.Request.Context(), id); err != nil {
		h.handleError(c, err, "delete_product")
		return
	}

	c.Status(http.StatusNoContent)
}

// CreateOrder
// @Summary      Create a new order
// @Description  Create a new order with multiple products
// @Tags         orders
// @Accept       json
// @Produce      json
// @Param        request body models.CreateOrderRequest true "Order details"
// @Success      201
// @Failure      400 {object} models.ErrorResponse
// @Failure      401 {object} models.ErrorResponse
// @Failure      404 {object} models.ErrorResponse
// @Failure      500 {object} models.ErrorResponse
// @Security     Bearer
// @Router       /orders [post]
func (h *handlerImpl) CreateOrder(c *gin.Context) {
	var req models.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.handleError(c, pkg.ErrInvalidInput, "create_order_validation")
		return
	}

	userID := c.GetString("user_id")
	if err := h.service.CreateOrder(c.Request.Context(), userID, &req); err != nil {
		h.handleError(c, err, "create_order")
		return
	}

	c.Status(http.StatusCreated)
}

// GetOrder
// @Summary      Get order by ID
// @Description  Get detailed information about a specific order
// @Tags         orders
// @Accept       json
// @Produce      json
// @Param        id path string true "Order ID"
// @Success      200 {object} models.Order
// @Failure      401 {object} models.ErrorResponse
// @Failure      404 {object} models.ErrorResponse
// @Failure      500 {object} models.ErrorResponse
// @Security     Bearer
// @Router       /orders/{id} [get]
func (h *handlerImpl) GetOrder(c *gin.Context) {
	id := c.Param("id")

	order, err := h.service.GetOrderByID(c.Request.Context(), id)
	if err != nil {
		h.handleError(c, err, "get_order")
		return
	}

	c.JSON(http.StatusOK, order)
}

// ListUserOrders
// @Summary      List user orders
// @Description  Get a list of all orders for the authenticated user
// @Tags         orders
// @Accept       json
// @Produce      json
// @Success      200 {array} models.Order
// @Failure      401 {object} models.ErrorResponse
// @Failure      500 {object} models.ErrorResponse
// @Security     Bearer
// @Router       /orders [get]
func (h *handlerImpl) ListUserOrders(c *gin.Context) {
	userID := c.GetString("user_id")

	orders, err := h.service.GetUserOrders(c.Request.Context(), userID)
	if err != nil {
		h.handleError(c, err, "list_user_orders")
		return
	}

	c.JSON(http.StatusOK, orders)
}

// UpdateOrderStatus
// @Summary      Update order status
// @Description  Update the status of an order (admin only)
// @Tags         orders
// @Accept       json
// @Produce      json
// @Param        id path string true "Order ID"
// @Param        request body models.UpdateOrderStatusRequest true "Order status update"
// @Success      200
// @Failure      400 {object} models.ErrorResponse
// @Failure      401 {object} models.ErrorResponse
// @Failure      403 {object} models.ErrorResponse
// @Failure      404 {object} models.ErrorResponse
// @Failure      500 {object} models.ErrorResponse
// @Security     Bearer
// @Router       /orders/{id}/status [put]
func (h *handlerImpl) UpdateOrderStatus(c *gin.Context) {
	id := c.Param("id")

	var req models.UpdateOrderStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.handleError(c, pkg.ErrInvalidInput, "update_order_status_validation")
		return
	}

	if err := h.service.UpdateStatus(c.Request.Context(), id, req.Status); err != nil {
		h.handleError(c, err, "update_order_status")
		return
	}

	c.Status(http.StatusOK)
}

// CancelOrder
// @Summary      Cancel order
// @Description  Cancel a pending order
// @Tags         orders
// @Accept       json
// @Produce      json
// @Param        id path string true "Order ID"
// @Success      200
// @Failure      400 {object} models.ErrorResponse "Order not in pending status"
// @Failure      401 {object} models.ErrorResponse
// @Failure      403 {object} models.ErrorResponse "Not the order owner"
// @Failure      404 {object} models.ErrorResponse
// @Failure      500 {object} models.ErrorResponse
// @Security     Bearer
// @Router       /orders/{id}/cancel [post]
func (h *handlerImpl) CancelOrder(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetString("user_id")

	if err := h.service.CancelOrder(c.Request.Context(), id, userID); err != nil {
		h.handleError(c, err, "cancel_order")
		return
	}

	c.Status(http.StatusOK)
}

// ErrorHandler provides centralized error handling with detailed logging and consistent responses
func (h *handlerImpl) handleError(ctx *gin.Context, err error, operation string) {
	// default error response
	statusCode := http.StatusInternalServerError
	errResp := models.ErrorResponse{
		Code:    "INTERNAL_ERROR",
		Message: "An unexpected error occurred",
	}

	// map specific errors to status codes and messages
	switch {
	case errors.Is(err, pkg.ErrInvalidCredentials):
		statusCode = http.StatusUnauthorized
		errResp.Code = "INVALID_CREDENTIALS"
		errResp.Message = "Invalid email or password"

	case errors.Is(err, pkg.ErrEmailTaken):
		statusCode = http.StatusConflict
		errResp.Code = "EMAIL_TAKEN"
		errResp.Message = "Email is already registered"

	case errors.Is(err, pkg.ErrNotFound):
		statusCode = http.StatusNotFound
		errResp.Code = "NOT_FOUND"
		errResp.Message = "Resource not found"

	case errors.Is(err, pkg.ErrUnauthorized):
		statusCode = http.StatusForbidden
		errResp.Code = "UNAUTHORIZED"
		errResp.Message = "You don't have permission to perform this action"

	case errors.Is(err, pkg.ErrInsufficientStock):
		statusCode = http.StatusBadRequest
		errResp.Code = "INSUFFICIENT_STOCK"
		errResp.Message = "One or more products are out of stock"

	case errors.Is(err, pkg.ErrOrderNotPending):
		statusCode = http.StatusBadRequest
		errResp.Code = "INVALID_ORDER_STATUS"
		errResp.Message = "Order cannot be modified in its current status"

	case errors.Is(err, pkg.ErrInvalidInput):
		statusCode = http.StatusBadRequest
		errResp.Code = "INVALID_INPUT"
		errResp.Message = "Invalid input provided"
		// add validation details if available
		if verr, ok := err.(interface{ Validation() map[string]string }); ok {
			errResp.Details = verr.Validation()
		}
	}

	// log error with context
	logAttrs := []slog.Attr{
		slog.String("operation", operation),
		slog.String("error_code", errResp.Code),
		slog.String("error", err.Error()),
		slog.Int("status_code", statusCode),
		slog.String("path", ctx.Request.URL.Path),
		slog.String("method", ctx.Request.Method),
	}

	// add user context if available
	if userID, exists := ctx.Get("user_id"); exists {
		logAttrs = append(logAttrs, slog.String("user_id", userID.(string)))
	}

	// log at appropriate level based on status code
	if statusCode >= 500 {
		slog.LogAttrs(ctx.Request.Context(), slog.LevelError, "Internal server error", logAttrs...)
	} else {
		slog.LogAttrs(ctx.Request.Context(), slog.LevelInfo, "Client error", logAttrs...)
	}

	ctx.JSON(statusCode, errResp)
}
