package main

import (
	"database/sql"
	"fmt"
	"time"
	_ "github.com/go-sql-driver/mysql"
	"platform/internal/config"
	"platform/internal/models"
	"platform/internal/repository/mysql"
	"platform/internal/repository/rabbitmq"
	"platform/pkg/logger"
	"strings"
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

	// Connect to MySQL with retries
	maxRetries := 10
	retryDelay := 5 * time.Second

	var db *sql.DB
	for i := 0; i < maxRetries; i++ {
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
			cfg.MySQL.User,
			cfg.MySQL.Password,
			cfg.MySQL.Host,
			cfg.MySQL.Port,
			cfg.MySQL.DBName,
		)

		db, err = sql.Open("mysql", dsn)
		if err != nil {
			logger.Error("Failed to open database connection", "attempt", i+1, "error", err)
			time.Sleep(retryDelay)
			continue
		}

		// Test database connection
		err = db.Ping()
		if err == nil {
			logger.Info("Successfully connected to database")
			break
		}

		logger.Error("Failed to ping database", "attempt", i+1, "error", err)
		time.Sleep(retryDelay)
	}

	if err != nil {
		logger.Fatal("Failed to connect to database after multiple attempts", err)
	}
	defer db.Close()

	// Initialize repositories
	requestRepo := mysql.NewRequestRepository(db)
	redirectRepo := mysql.NewRedirectRepository(db)
	consumer, err := rabbitmq.NewConsumer(cfg)
	if err != nil {
		logger.Fatal("Failed to initialize RabbitMQ consumer", err)
	}
	defer consumer.Close()

	// Start consuming messages
	logger.Info("Starting database worker")
	if err := consumer.Consume(func(request *models.Request) error {
		// Save request to database
		if err := requestRepo.SaveRequest(request); err != nil {
			return fmt.Errorf("failed to save request: %w", err)
		}

		// Extract hash from request URL
		hash := extractHashFromURL(request.RequestURL)
		if hash != "" {
			// Get redirect URL
			redirectURL, err := redirectRepo.GetRedirectURL(hash)
			if err != nil {
				return fmt.Errorf("failed to get redirect URL: %w", err)
			}

			if redirectURL != "" {
				// Create redirect record
				redirect := &models.Redirect{
					RequestLogID:     request.ID,
					OriginalURL:      request.RequestURL,
					RedirectURL:      redirectURL,
					RedirectType:     determineRedirectType(redirectURL),
					RedirectStatus:   302, // Temporary redirect
					RedirectTimestamp: request.Timestamp,
				}

				// Save redirect record
				if err := redirectRepo.SaveRedirect(redirect); err != nil {
					return fmt.Errorf("failed to save redirect: %w", err)
				}
			}
		}

		return nil
	}); err != nil {
		logger.Fatal("Failed to start consuming messages", err)
	}

	// Keep the worker running
	select {}
}

func extractHashFromURL(url string) string {
	// Extract hash from URL path
	// Example: /345678?click_id=clic-hash -> 345678
	parts := strings.Split(url, "/")
	if len(parts) >= 2 {
		// Remove any query parameters
		hash := strings.Split(parts[1], "?")[0]
		return hash
	}
	return ""
}

func determineRedirectType(redirectURL string) string {
	// Extract the domain from the URL
	// Example: http://site1.com/special-offer -> site1
	parts := strings.Split(redirectURL, "/")
	if len(parts) >= 3 {
		domain := strings.Split(parts[2], ".")[0]
		return domain
	}
	return "unknown"
} 