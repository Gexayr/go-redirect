package api

import (
	"github.com/gin-gonic/gin"
	"platform/internal/api/handlers"
	"platform/internal/api/middleware"
	"platform/internal/config"
	"platform/internal/database"
	"platform/internal/repository/mysql"
	"platform/internal/repository/rabbitmq"
	"platform/pkg/logger"
)

func SetupRouter(publisher *rabbitmq.Publisher) *gin.Engine {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Fatal("Failed to load configuration", err)
	}

	// Initialize database connection
	if err := database.InitDB(cfg); err != nil {
		logger.Fatal("Failed to initialize database", err)
	}

	// Initialize repositories
	redirectRepo := mysql.NewRedirectRepository(database.GetDB())

	// Initialize request handler
	requestHandler := handlers.NewRequestHandler(publisher, redirectRepo)

	// Create router
	router := gin.New()

	// Add middleware
	router.Use(middleware.RequestID())
	router.Use(middleware.Logging())

	// Health check
	router.GET("/health", handlers.HealthCheck)

	// Hash endpoint with dynamic hash parameter
	router.GET("/:hash", requestHandler.ProcessRequest)

	return router
} 