package main

import (
	"database/sql"
	"fmt"
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

	// Connect to MySQL
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		cfg.MySQL.User,
		cfg.MySQL.Password,
		cfg.MySQL.Host,
		cfg.MySQL.Port,
		cfg.MySQL.DBName,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		logger.Fatal("Failed to connect to MySQL", err)
	}
	defer db.Close()

	// Test database connection
	if err := db.Ping(); err != nil {
		logger.Fatal("Failed to ping database", err)
	}

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
	// Example: http://localhost:8080/123456?click_id=click-hash -> 123456
	parts := strings.Split(url, "/")
	if len(parts) >= 3 {
		return parts[2]
	}
	return ""
}

func determineRedirectType(redirectURL string) string {
	switch {
	case strings.Contains(redirectURL, "site1.com"):
		return string(models.RedirectTypeSite1)
	case strings.Contains(redirectURL, "site2.com"):
		return string(models.RedirectTypeSite2)
	case strings.Contains(redirectURL, "site3.com"):
		return string(models.RedirectTypeSite3)
	default:
		return ""
	}
} 