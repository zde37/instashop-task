package handler

import "github.com/gin-gonic/gin"

type Handler interface {
	Register(ctx *gin.Context)
	Login(ctx *gin.Context)
	Logout(ctx *gin.Context)

	CreateProduct(ctx *gin.Context)
	GetProduct(ctx *gin.Context)
	ListProducts(ctx *gin.Context)
	UpdateProduct(ctx *gin.Context)
	DeleteProduct(ctx *gin.Context)

	CreateOrder(ctx *gin.Context)
	GetOrder(ctx *gin.Context)
	ListUserOrders(ctx *gin.Context)
	UpdateOrderStatus(ctx *gin.Context)
	CancelOrder(ctx *gin.Context)
}
