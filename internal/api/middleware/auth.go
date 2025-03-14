package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"platform/internal/auth"
	"platform/internal/repository/mysql"
	"platform/pkg/logger"
	"strings"
)

func Auth(clientRepo *mysql.ClientRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		// Extract token from Bearer scheme
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		// Validate JWT token
		claims, err := auth.ValidateToken(parts[1])
		if err != nil {
			logger.Error("Failed to validate token", "error", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Get client by ID
		client, err := clientRepo.GetByID(claims.ClientID)
		if err != nil {
			logger.Error("Failed to get client", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to authenticate"})
			c.Abort()
			return
		}

		if client == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Set client ID in context
		c.Set("client_id", client.ID)
		c.Next()
	}
} 