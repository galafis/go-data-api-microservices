package middleware

import (
	"net/http"
	"strings"

	"github.com/galafis/go-data-api-microservices/internal/auth"
	"github.com/galafis/go-data-api-microservices/internal/models"
	"github.com/galafis/go-data-api-microservices/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// AuthMiddleware represents the authentication middleware
type AuthMiddleware struct {
	jwtService *auth.JWTService
}

// NewAuthMiddleware creates a new authentication middleware
func NewAuthMiddleware(jwtService *auth.JWTService) *AuthMiddleware {
	return &AuthMiddleware{
		jwtService: jwtService,
	}
}

// AuthRequired is a middleware that checks if the user is authenticated
func (m *AuthMiddleware) AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header is required"})
			c.Abort()
			return
		}

		// Check if the header has the Bearer prefix
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header format"})
			c.Abort()
			return
		}

		// Extract and validate the token
		tokenString := parts[1]
		claims, err := m.jwtService.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			c.Abort()
			return
		}

		// Check if it's an access token
		if claims.TokenType != "access" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token type"})
			c.Abort()
			return
		}

		// Parse user ID
		userID, err := uuid.Parse(claims.UserID)
		if err != nil {
			logger.Errorf("Invalid user ID in token: %v", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		// Set user information in the context
		c.Set("user_id", userID)
		c.Set("email", claims.Email)
		c.Set("role", claims.Role)

		c.Next()
	}
}

// RoleRequired is a middleware that checks if the user has the required role
func (m *AuthMiddleware) RoleRequired(roles ...models.Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user role from context
		role, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		userRole, ok := role.(models.Role)
		if !ok {
			logger.Error("Failed to cast role to models.Role")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			c.Abort()
			return
		}

		// Check if user has one of the required roles
		hasRole := false
		for _, r := range roles {
			if userRole == r {
				hasRole = true
				break
			}
		}

		if !hasRole {
			c.JSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// AuthRequired is a shorthand function for the auth middleware
func AuthRequired() gin.HandlerFunc {
	// This is a placeholder that should be replaced with a proper implementation
	// that gets the JWT service from the application context
	jwtService := auth.NewJWTService(nil) // This should be properly initialized
	middleware := NewAuthMiddleware(jwtService)
	return middleware.AuthRequired()
}

// RoleRequired is a shorthand function for the role middleware
func RoleRequired(roles ...models.Role) gin.HandlerFunc {
	// This is a placeholder that should be replaced with a proper implementation
	// that gets the JWT service from the application context
	jwtService := auth.NewJWTService(nil) // This should be properly initialized
	middleware := NewAuthMiddleware(jwtService)
	return middleware.RoleRequired(roles...)
}

