package handlers

import (
	"net/http"
	"time"

	"github.com/galafis/go-data-api-microservices/internal/auth"
	"github.com/galafis/go-data-api-microservices/internal/models"
	"github.com/galafis/go-data-api-microservices/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// AuthHandler handles authentication requests
type AuthHandler struct {
	jwtService      auth.JWTService
	passwordService auth.PasswordService
	userRepository  UserRepository
}

// UserRepository defines the interface for user operations
type UserRepository interface {
	FindByEmail(email string) (*models.User, error)
	FindByID(userID uuid.UUID) (*models.User, error)
	Create(user *models.User) error
	Update(user *models.User) error
	Delete(userID uuid.UUID) error
	UpdateRefreshToken(userID uuid.UUID, token string) error
	ClearRefreshToken(userID uuid.UUID) error
}

// NewAuthHandler creates a new authentication handler
func NewAuthHandler(jwtService auth.JWTService, passwordService auth.PasswordService, userRepository UserRepository) *AuthHandler {
	return &AuthHandler{
		jwtService:      jwtService,
		passwordService: passwordService,
		userRepository:  userRepository,
	}
}

// Register handles user registration
// @Summary Register a new user
// @Description Register a new user with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.RegisterRequest true "Registration request"
// @Success 201 {object} models.AuthResponse "User registered successfully"
// @Failure 400 {object} ErrorResponse "Invalid request"
// @Failure 409 {object} ErrorResponse "Email already exists"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if email already exists
	existingUser, err := h.userRepository.FindByEmail(req.Email)
	if err != nil {
		logger.Errorf("Error checking existing user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	if existingUser != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
		return
	}

	// Check password strength
	if strong, reason := h.passwordService.IsStrongPassword(req.Password); !strong {
		c.JSON(http.StatusBadRequest, gin.H{"error": reason})
		return
	}

	// Hash password
	hashedPassword, err := h.passwordService.HashPassword(req.Password)
	if err != nil {
		logger.Errorf("Error hashing password: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	// Create user
	now := time.Now()
	user := &models.User{
		ID:        uuid.New(),
		Email:     req.Email,
		Password:  hashedPassword,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Role:      models.RoleUser,
		Active:    true,
		Verified:  false,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := h.userRepository.Create(user); err != nil {
		logger.Errorf("Error creating user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	// Generate tokens
	accessToken, err := h.jwtService.GenerateAccessToken(user)
	if err != nil {
		logger.Errorf("Error generating access token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	refreshToken, err := h.jwtService.GenerateRefreshToken(user)
	if err != nil {
		logger.Errorf("Error generating refresh token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	// Store refresh token
	if err := h.userRepository.UpdateRefreshToken(user.ID, refreshToken); err != nil {
		logger.Errorf("Error storing refresh token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	// Return response
	c.JSON(http.StatusCreated, models.AuthResponse{
		User: models.UserResponse{
			ID:        user.ID,
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Role:      user.Role,
			Active:    user.Active,
			Verified:  user.Verified,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
		AccessToken:  accessToken,
		TokenType:    "Bearer",
		ExpiresIn:    int64(h.jwtService.GetTokenExpiry("access").Seconds()),
		RefreshToken: refreshToken,
	})
}

// Login handles user login
// @Summary Login a user
// @Description Login a user with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.LoginRequest true "Login request"
// @Success 200 {object} models.AuthResponse "User logged in successfully"
// @Failure 400 {object} ErrorResponse "Invalid request"
// @Failure 401 {object} ErrorResponse "Invalid credentials"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find user by email
	user, err := h.userRepository.FindByEmail(req.Email)
	if err != nil {
		logger.Errorf("Error finding user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Check if user is active
	if !user.Active {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Account is inactive"})
		return
	}

	// Verify password
	if err := h.passwordService.VerifyPassword(user.Password, req.Password); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Update last login time
	now := time.Now()
	user.LastLoginAt = &now
	user.UpdatedAt = now

	// Generate tokens
	accessToken, err := h.jwtService.GenerateAccessToken(user)
	if err != nil {
		logger.Errorf("Error generating access token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	refreshToken, err := h.jwtService.GenerateRefreshToken(user)
	if err != nil {
		logger.Errorf("Error generating refresh token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	// Store refresh token
	if err := h.userRepository.UpdateRefreshToken(user.ID, refreshToken); err != nil {
		logger.Errorf("Error storing refresh token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	// Return response
	c.JSON(http.StatusOK, models.AuthResponse{
		User: models.UserResponse{
			ID:        user.ID,
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Role:      user.Role,
			Active:    user.Active,
			Verified:  user.Verified,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
		AccessToken:  accessToken,
		TokenType:    "Bearer",
		ExpiresIn:    int64(h.jwtService.GetTokenExpiry("access").Seconds()),
		RefreshToken: refreshToken,
	})
}

// RefreshToken handles token refresh
// @Summary Refresh access token
// @Description Refresh access token using refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.RefreshTokenRequest true "Refresh token request"
// @Success 200 {object} models.AuthResponse "Token refreshed successfully"
// @Failure 400 {object} ErrorResponse "Invalid request"
// @Failure 401 {object} ErrorResponse "Invalid refresh token"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req models.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate refresh token
	claims, err := h.jwtService.ValidateToken(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}

	// Check token type
	if claims.TokenType != "refresh" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token type"})
		return
	}

	// Parse user ID
	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		logger.Errorf("Invalid user ID in token: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	// Find user by ID
	user, err := h.userRepository.FindByID(userID)
	if err != nil {
		logger.Errorf("Error finding user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	// Check if user is active
	if !user.Active {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Account is inactive"})
		return
	}

	// Check if refresh token matches
	if user.RefreshToken != req.RefreshToken {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}

	// Generate new tokens
	accessToken, err := h.jwtService.GenerateAccessToken(user)
	if err != nil {
		logger.Errorf("Error generating access token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	refreshToken, err := h.jwtService.GenerateRefreshToken(user)
	if err != nil {
		logger.Errorf("Error generating refresh token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	// Store new refresh token
	if err := h.userRepository.UpdateRefreshToken(user.ID, refreshToken); err != nil {
		logger.Errorf("Error storing refresh token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	// Return response
	c.JSON(http.StatusOK, models.AuthResponse{
		User: models.UserResponse{
			ID:        user.ID,
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Role:      user.Role,
			Active:    user.Active,
			Verified:  user.Verified,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
		AccessToken:  accessToken,
		TokenType:    "Bearer",
		ExpiresIn:    int64(h.jwtService.GetTokenExpiry("access").Seconds()),
		RefreshToken: refreshToken,
	})
}

// Logout handles user logout
// @Summary Logout a user
// @Description Logout a user and invalidate refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} SuccessResponse "User logged out successfully"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Clear refresh token
	if err := h.userRepository.ClearRefreshToken(userID.(uuid.UUID)); err != nil {
		logger.Errorf("Error clearing refresh token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

// Register is a placeholder handler for user registration
func Register(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Register endpoint"})
}

// Login is a placeholder handler for user login
func Login(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Login endpoint"})
}

// RefreshToken is a placeholder handler for token refresh
func RefreshToken(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Refresh token endpoint"})
}

// Logout is a placeholder handler for user logout
func Logout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Logout endpoint"})
}

