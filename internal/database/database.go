package database

import (
	"database/sql"
	"fmt"
	"time"
	_ "github.com/go-sql-driver/mysql"
	"platform/internal/config"
	"platform/pkg/logger"
)

var db *sql.DB

// InitDB initializes the database connection with retries
func InitDB(cfg *config.Config) error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		cfg.MySQL.User,
		cfg.MySQL.Password,
		cfg.MySQL.Host,
		cfg.MySQL.Port,
		cfg.MySQL.DBName,
	)

	maxRetries := 10
	retryDelay := 5 * time.Second

	var err error
	for i := 0; i < maxRetries; i++ {
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
			return nil
		}

		logger.Error("Failed to ping database", "attempt", i+1, "error", err)
		time.Sleep(retryDelay)
	}

	return fmt.Errorf("failed to connect to database after %d attempts: %w", maxRetries, err)
}

// GetDB returns the database connection
func GetDB() *sql.DB {
	return db
}

// CloseDB closes the database connection
func CloseDB() {
	if db != nil {
		db.Close()
	}
} 