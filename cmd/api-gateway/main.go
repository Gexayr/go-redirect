package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"platform/internal/api"
	"platform/internal/config"
	"platform/internal/database"
	"platform/internal/repository/rabbitmq"
	"platform/pkg/logger"
)

func main() {
	// Initialize logger
	if err := logger.Init(); err != nil {
		logger.Fatal("Failed to initialize logger", err)
	}

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Fatal("Failed to load configuration", err)
	}

	// Initialize RabbitMQ publisher
	publisher, err := rabbitmq.NewPublisher(cfg)
	if err != nil {
		logger.Fatal("Failed to initialize RabbitMQ publisher", err)
	}
	defer publisher.Close()

	// Setup router
	router := api.SetupRouter(publisher)

	// Create a channel to listen for shutdown signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Start server in a goroutine
	go func() {
		addr := fmt.Sprintf(":%s", cfg.Server.Port)
		logger.Info("Starting server", "port", cfg.Server.Port)
		if err := router.Run(addr); err != nil {
			logger.Error("Failed to start server", "error", err)
		}
	}()

	// Wait for shutdown signal
	<-quit

	// Graceful shutdown
	logger.Info("Shutting down server...")
	database.CloseDB()
	logger.Info("Server stopped")
} 