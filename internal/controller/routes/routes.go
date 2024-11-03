package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zde37/instashop-task/internal/controller/handler"
	"github.com/zde37/instashop-task/internal/controller/middlewares"
	"github.com/zde37/instashop-task/pkg"
)

// RegisterAllRoutes registers all the routes for the application.
func RegisterAllRoutes(rg *gin.RouterGroup, handler handler.Handler, jwt *pkg.JWTMaker) {
	// health check
	rg.GET("/ping", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "pong")
	})

	rg.POST("/register", handler.Register)
	rg.POST("/login", handler.Login)

	// protected routes
	api := rg.Group("/")
	api.Use(middlewares.Auth(jwt))
	{
		products := api.Group("/products")
		{
			products.GET("", handler.ListProducts)
			products.GET("/:id", handler.GetProduct)

			products.Use(middlewares.AdminRequired())
			{
				products.POST("", handler.CreateProduct)
				products.PUT("/:id", handler.UpdateProduct)
				products.DELETE("/:id", handler.DeleteProduct)
			}
		}

		orders := api.Group("/orders")
		{
			orders.POST("", handler.CreateOrder)
			orders.GET("", handler.ListUserOrders)
			orders.GET("/:id", handler.GetOrder)
			orders.POST("/:id/cancel", handler.CancelOrder)

			orders.Use(middlewares.AdminRequired())
			{
				orders.PUT("/:id/status", handler.UpdateOrderStatus)
			}
		}
	}

	api.POST("/logout", handler.Logout)
}
