package controller

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "github.com/zde37/instashop-task/docs"
	"github.com/zde37/instashop-task/internal/config"
	"github.com/zde37/instashop-task/internal/controller/handler"
	"github.com/zde37/instashop-task/internal/controller/routes"
	"github.com/zde37/instashop-task/internal/repository"
	"github.com/zde37/instashop-task/internal/service"
	"github.com/zde37/instashop-task/pkg"
)

// Controller represents the main application structure
type Controller struct {
	router     *gin.Engine
	config     *config.Config
	db         *pgxpool.Pool
	handler    handler.Handler
	httpServer *http.Server
}

// New creates a new instance of Controller
func New() *Controller {
	return &Controller{
		router: gin.Default(),
	}
}

// bootstrap initializes the Controller and its dependencies
func (c *Controller) bootstrap() error {
	var err error

	// load config
	c.config, err = config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %v", err)
	}

	ctx := context.Background()
	c.db, err = config.SetupDatabase(ctx, c.config.DSN, "file://migrations")
	if err != nil {
		return fmt.Errorf("failed to setup database: %v", err)
	}

	// initialize jwt
	jwtMaker, err := pkg.NewJWTMaker(c.config.JWTSecretKey)
	if err != nil {
		return fmt.Errorf("failed to initialize jwt maker: %v", err)
	}

	// initialize repository, service, and handlers
	repo := repository.New(c.db)
	srvc := service.New(repo, jwtMaker)
	c.handler = handler.New(srvc)

	if c.config.Environment == pkg.Production {
		gin.SetMode(gin.ReleaseMode)
	}

	// register password validator
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		if err := v.RegisterValidation("validpassword", pkg.ValidatePassword); err != nil {
			slog.Error(err.Error())
		}
	}

	c.setupRoutes(jwtMaker)
	c.configureHTTPServer()
	return nil
}

func (c *Controller) setupRoutes(jwtMaker *pkg.JWTMaker) {
	v1RouteGroup := c.router.Group("/api/v1")
	v1RouteGroup.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	routes.RegisterAllRoutes(v1RouteGroup, c.handler, jwtMaker)
}

func (c *Controller) configureHTTPServer() {
	c.httpServer = &http.Server{
		Addr:         fmt.Sprintf(":%s", c.config.Port),
		Handler:      c.router,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
}

// Serve starts the Controller server
func (c *Controller) Serve() error {
	if err := c.bootstrap(); err != nil {
		return err
	}

	// start the server
	go func() {
		slog.Info("starting server...", slog.String("Port", c.config.Port))
		if err := c.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("failed to start server", slog.String("err", err.Error()))
		}
	}()

	return c.gracefulShutdown()
}

// gracefulShutdown handles graceful shutdown of all services
func (c *Controller) gracefulShutdown() error {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("shutting down services...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	c.httpServer.SetKeepAlivesEnabled(false)

	if err := c.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("error shutting down server: %v", err)
	}

	c.db.Close()
	slog.Info("services stopped gracefully")
	return nil
}
