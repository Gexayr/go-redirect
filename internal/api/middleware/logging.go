package middleware

import (
	"github.com/gin-gonic/gin"
	"platform/pkg/logger"
	"time"
)

func Logging() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		// Log request details
		logger.Info("Request processed",
			"method", c.Request.Method,
			"path", path,
			"status", c.Writer.Status(),
			"latency", time.Since(start),
			"client_ip", c.ClientIP(),
			"query", query,
		)
	}
} 