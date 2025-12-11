package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/yourusername/go-api-service/internal/config"
	"github.com/yourusername/go-api-service/internal/handler"
	"github.com/yourusername/go-api-service/internal/repository"
	"github.com/yourusername/go-api-service/internal/service"
	"github.com/yourusername/go-api-service/pkg/database"
	"github.com/yourusername/go-api-service/pkg/logger"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize logger
	appLogger := logger.New(cfg.Log.Level)
	appLogger.Info("Starting application", "environment", cfg.Server.Environment)

	// Initialize database connection
	dbConfig := database.DBConfig{
		DSN:             cfg.GetDSN(),
		MaxOpenConns:    25,
		MaxIdleConns:    5,
		ConnMaxLifetime: 5 * time.Minute,
	}

	db, err := database.NewPostgresConnection(dbConfig)
	if err != nil {
		appLogger.Error("Failed to connect to database", "error", err)
		log.Fatalf("Database connection error: %v", err)
	}
	defer db.Close()
	appLogger.Info("Database connection established")

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)

	// Initialize services
	userService := service.NewUserService(userRepo, appLogger)

	// Initialize handlers
	userHandler := handler.NewUserHandler(userService, appLogger)
	healthHandler := handler.NewHealthHandler()

	// Setup router
	router := handler.NewRouter(userHandler, healthHandler, appLogger)

	// Create HTTP server
	serverAddr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	srv := &http.Server{
		Addr:         serverAddr,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		appLogger.Info("Server starting", "address", serverAddr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			appLogger.Error("Server failed to start", "error", err)
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	appLogger.Info("Server shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		appLogger.Error("Server forced to shutdown", "error", err)
		log.Fatal(err)
	}

	appLogger.Info("Server stopped gracefully")
}
