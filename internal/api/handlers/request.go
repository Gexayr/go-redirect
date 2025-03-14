package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"platform/internal/models"
	"platform/internal/repository/mysql"
	"platform/internal/repository/rabbitmq"
	"platform/pkg/logger"
	"time"
)

type RequestHandler struct {
	publisher        *rabbitmq.Publisher
	redirectRepo     *mysql.RedirectRepository
}

func NewRequestHandler(publisher *rabbitmq.Publisher, redirectRepo *mysql.RedirectRepository) *RequestHandler {
	return &RequestHandler{
		publisher:    publisher,
		redirectRepo: redirectRepo,
	}
}

func (h *RequestHandler) ProcessRequest(c *gin.Context) {
	// Only accept GET requests
	if c.Request.Method != http.MethodGet {
		logger.Error("Invalid request method", "method", c.Request.Method)
		c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "Only GET requests are allowed"})
		return
	}

	// Get hash from URL path
	hash := c.Param("hash")
	if hash == "" {
		logger.Error("Missing hash in URL path")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing hash in URL path"})
		return
	}

	// Get click_id parameter
	clickID := c.Query("click_id")
	if clickID == "" {
		logger.Error("Missing click_id parameter")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing click_id parameter"})
		return
	}

	// Convert headers to JSON
	headers := make(map[string]string)
	for k, v := range c.Request.Header {
		if len(v) > 0 {
			headers[k] = v[0]
		}
	}
	headersJSON, err := json.Marshal(headers)
	if err != nil {
		headersJSON = []byte("{}")
	}

	// Create request log
	request := &models.Request{
		Timestamp:      time.Now(),
		IPAddress:      c.ClientIP(),
		RequestURL:     c.Request.URL.String(),
		RequestMethod:  c.Request.Method,
		RequestHeaders: headersJSON,
	}

	// Get redirect URL from database
	redirectURL, err := h.redirectRepo.GetRedirectURL(hash)
	if err != nil {
		logger.Error("Failed to get redirect URL", "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process request"})
		return
	}

	if redirectURL != "" {
		// Store request in RabbitMQ
		if err := h.publisher.PublishRequest(request); err != nil {
			logger.Error("Failed to publish request", "error", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process request"})
			return
		}

		// Append click_id parameter to redirect URL
		finalURL := fmt.Sprintf("%s?click_id=%s", redirectURL, clickID)

		// Redirect to the appropriate site
		c.Redirect(http.StatusTemporaryRedirect, finalURL)
		return
	}

	// If no redirect URL found, just store the request
	if err := h.publisher.PublishRequest(request); err != nil {
		logger.Error("Failed to publish request", "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process request"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Request processed successfully",
		"request": request,
	})
} 