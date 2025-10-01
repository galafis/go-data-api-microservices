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

// UserHandler handles user operations
type UserHandler struct {
	userRepository  UserRepository
	passwordService auth.PasswordService
}

// UserRepository interface is defined in auth.go

// NewUserHandler creates a new user handler
func NewUserHandler(userRepository UserRepository, passwordService auth.PasswordService) *UserHandler {
	return &UserHandler{
		userRepository:  userRepository,
		passwordService: passwordService,
	}
}

// GetCurrentUser handles getting the current user
// @Summary Get current user
// @Description Get the current authenticated user
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.UserResponse "User retrieved successfully"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 404 {object} ErrorResponse "User not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /users/me [get]
func (h *UserHandler) GetCurrentUser(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Get user
	user, err := h.userRepository.FindByID(userID.(uuid.UUID))
	if err != nil {
		logger.Errorf("Error finding user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Convert to response
	response := models.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Role:      user.Role,
		Active:    user.Active,
		Verified:  user.Verified,
		Metadata:  user.Metadata,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	c.JSON(http.StatusOK, response)
}

// UpdateCurrentUser handles updating the current user
// @Summary Update current user
// @Description Update the current authenticated user
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.UpdateUserRequest true "User update request"
// @Success 200 {object} models.UserResponse "User updated successfully"
// @Failure 400 {object} ErrorResponse "Invalid request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 404 {object} ErrorResponse "User not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /users/me [put]
func (h *UserHandler) UpdateCurrentUser(c *gin.Context) {
	// Parse request
	var req models.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Get user
	user, err := h.userRepository.FindByID(userID.(uuid.UUID))
	if err != nil {
		logger.Errorf("Error finding user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Update user
	if req.FirstName != "" {
		user.FirstName = req.FirstName
	}
	if req.LastName != "" {
		user.LastName = req.LastName
	}
	if req.Password != "" {
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
		user.Password = hashedPassword
	}
	if req.Metadata != nil {
		if user.Metadata == nil {
			user.Metadata = req.Metadata
		} else {
			// Merge metadata
			for k, v := range req.Metadata {
				user.Metadata[k] = v
			}
		}
	}
	user.UpdatedAt = time.Now()

	// Save user
	if err := h.userRepository.Update(user); err != nil {
		logger.Errorf("Error updating user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	// Convert to response
	response := models.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Role:      user.Role,
		Active:    user.Active,
		Verified:  user.Verified,
		Metadata:  user.Metadata,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	c.JSON(http.StatusOK, response)
}

// DeleteCurrentUser handles deleting the current user
// @Summary Delete current user
// @Description Delete the current authenticated user
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 204 "User deleted successfully"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 404 {object} ErrorResponse "User not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /users/me [delete]
func (h *UserHandler) DeleteCurrentUser(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Check if user exists
	user, err := h.userRepository.FindByID(userID.(uuid.UUID))
	if err != nil {
		logger.Errorf("Error finding user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Delete user
	if err := h.userRepository.Delete(userID.(uuid.UUID)); err != nil {
		logger.Errorf("Error deleting user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.Status(http.StatusNoContent)
}

// GetCurrentUser is a placeholder handler for getting the current user
func GetCurrentUser(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Get current user endpoint"})
}

// UpdateCurrentUser is a placeholder handler for updating the current user
func UpdateCurrentUser(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Update current user endpoint"})
}

// DeleteCurrentUser is a placeholder handler for deleting the current user
func DeleteCurrentUser(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Delete current user endpoint"})
}

