package handlers

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"platform/internal/auth"
	"platform/internal/models"
	"platform/internal/repository/mysql"
	"platform/pkg/logger"
)

type ClientHandler struct {
	clientRepo   *mysql.ClientRepository
	redirectRepo *mysql.RedirectRepository
}

func NewClientHandler(clientRepo *mysql.ClientRepository, redirectRepo *mysql.RedirectRepository) *ClientHandler {
	return &ClientHandler{
		clientRepo:   clientRepo,
		redirectRepo: redirectRepo,
	}
}

// Register handles client registration
func (h *ClientHandler) Register(c *gin.Context) {
	var registration models.ClientRegistration
	if err := c.ShouldBindJSON(&registration); err != nil {
		logger.Error("Invalid registration data", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid registration data"})
		return
	}

	// Check if username already exists
	existingClient, err := h.clientRepo.GetByUsername(registration.Username)
	if err != nil {
		logger.Error("Failed to check username", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process registration"})
		return
	}
	if existingClient != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Username already exists"})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registration.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Error("Failed to hash password", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process registration"})
		return
	}

	// Create client
	client := &models.Client{
		Username:     registration.Username,
		PasswordHash: string(hashedPassword),
		Email:        registration.Email,
	}

	if err := h.clientRepo.Create(client); err != nil {
		logger.Error("Failed to create client", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process registration"})
		return
	}

	// Generate JWT token
	token, err := auth.GenerateToken(client.ID, client.Username)
	if err != nil {
		logger.Error("Failed to generate token", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process registration"})
		return
	}

	// Return client response with token
	response := &models.ClientResponse{
		ID:        client.ID,
		Username:  client.Username,
		Email:     client.Email,
		CreatedAt: client.CreatedAt,
		Token:     token,
	}

	c.JSON(http.StatusCreated, response)
}

// Login handles client login
func (h *ClientHandler) Login(c *gin.Context) {
	var login models.ClientLogin
	if err := c.ShouldBindJSON(&login); err != nil {
		logger.Error("Invalid login data", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid login data"})
		return
	}

	// Get client by username
	client, err := h.clientRepo.GetByUsername(login.Username)
	if err != nil {
		logger.Error("Failed to get client", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process login"})
		return
	}
	if client == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Validate password
	if !h.clientRepo.ValidatePassword(client, login.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate JWT token
	token, err := auth.GenerateToken(client.ID, client.Username)
	if err != nil {
		logger.Error("Failed to generate token", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process login"})
		return
	}

	// Return client response with token
	response := &models.ClientResponse{
		ID:        client.ID,
		Username:  client.Username,
		Email:     client.Email,
		CreatedAt: client.CreatedAt,
		Token:     token,
	}

	c.JSON(http.StatusOK, response)
}

// CreateRedirectMapping creates a new redirect mapping for the authenticated client
func (h *ClientHandler) CreateRedirectMapping(c *gin.Context) {
	clientID, exists := c.Get("client_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var mapping models.RedirectMappingCreate
	if err := c.ShouldBindJSON(&mapping); err != nil {
		logger.Error("Invalid redirect mapping data", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid redirect mapping data"})
		return
	}

	redirectMapping := &models.RedirectMapping{
		RedirectURL:     mapping.RedirectURL,
		RedirectURLBlack: mapping.RedirectURLBlack,
	}

	if err := h.redirectRepo.CreateRedirectMapping(clientID.(int64), redirectMapping); err != nil {
		logger.Error("Failed to create redirect mapping", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create redirect mapping"})
		return
	}

	c.JSON(http.StatusCreated, redirectMapping)
}

// GetRedirectMappings returns all redirect mappings for the authenticated client
func (h *ClientHandler) GetRedirectMappings(c *gin.Context) {
	clientID, exists := c.Get("client_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	mappings, err := h.redirectRepo.GetClientRedirectMappings(clientID.(int64))
	if err != nil {
		logger.Error("Failed to get redirect mappings", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get redirect mappings"})
		return
	}

	c.JSON(http.StatusOK, mappings)
} 