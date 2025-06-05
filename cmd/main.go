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

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
)

func main() {
	// Load .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("No .env file found — relying on system environment variables")
	}
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
	log.Printf("Loaded config: %+v\n", cfg)

	// Initialize logging
	logger, err := logging.NewRotateLogger(cfg.Logging.Directory, cfg.Logging.Filename)
	log.Printf("logger initialised succesfully")
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Close()

	// Database connection
	database, err := db.ConnectMySQL(cfg.Database)
	log.Println("Database ziko")
	if err != nil {
		log.Printf("cannot load db")
		logger.Fatal().Err(err).Msg("Failed to connect to database")
	}
	defer database.Close()
	logger.Info().Msg("Successfully connected to database")

	// Run migrations
	log.Println("About to create tables...")
	if err := db.RunSQLFile(database, "db.sql"); err != nil {
		logger.Fatal().Err(err).Msg("Failed to apply DB schema")
	}

	// Initialize services
	authService := auth.NewOpenIDService(cfg.Auth)
	notificationService := services.NewNotificationService(cfg.Notifications)

	orderService := services.NewOrderService(
		repositories.NewOrderRepository(database),
		*repositories.NewProductRepository(database),
		*repositories.NewCustomerRepository(database),
		*notificationService,
	)

	productService := services.NewProductService(
		*repositories.NewProductRepository(database),
		*repositories.NewCategoryRepository(database),
	)

	categoryService := services.NewCategoryService(
		*repositories.NewCategoryRepository(database),
	)

	// Initialize HTTP handlers
	handler := handlers.NewAPIHandler(
		authService,
		orderService,
		productService,
		categoryService,
		&zerolog.Logger{},
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
