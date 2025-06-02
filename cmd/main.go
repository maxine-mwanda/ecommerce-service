// cmd/main.go
package main

import (
	"context"
	"ecommerce-service/internal/auth"
	"ecommerce-service/internal/config"
	"ecommerce-service/internal/db"
	"ecommerce-service/internal/handlers"
	"ecommerce-service/internal/logging"
	"ecommerce-service/internal/repositories"
	"ecommerce-service/internal/services"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize logging
	logger, err := logging.NewRotateLogger(cfg.Logging.Directory, cfg.Logging.Filename)
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Close()

	// Database connection
	database, err := db.ConnectMySQL(cfg.Database)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to connect to database")
	}
	defer database.Close()

	// Run migrations
	if err := db.RunMigrations(database, cfg.Database); err != nil {
		logger.Fatal().Err(err).Msg("Failed to run migrations")
	}

	// Initialize repositories
	orderRepo := repositories.NewOrderRepository(database)
	productRepo := repositories.NewProductRepository(database)
	categoryRepo := repositories.NewCategoryRepository(database)
	customerRepo := repositories.NewCustomerRepository(database)

	// Initialize services
	authService := auth.NewOpenIDService(cfg.Auth)
	notificationService := services.NewNotificationService(cfg.Notifications)
	orderService := services.NewOrderService(orderRepo, productRepo, customerRepo, notificationService)
	productService := services.NewProductService(productRepo, categoryRepo)
	categoryService := services.NewCategoryService(categoryRepo)

	// Initialize HTTP handlers
	handler := handlers.NewAPIHandler(
		authService,
		orderService,
		productService,
		categoryService,
		logger,
	)

	// Create HTTP server
	server := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: handler,
	}

	// Start server in a goroutine
	go func() {
		logger.Info().Str("port", cfg.Server.Port).Msg("Starting server")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal().Err(err).Msg("Server failed")
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info().Msg("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error().Err(err).Msg("Server forced to shutdown")
	}

	logger.Info().Msg("Server exiting")
}
