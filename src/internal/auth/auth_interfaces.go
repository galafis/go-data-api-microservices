package auth

import (
	"time"

	"github.com/galafis/go-data-api-microservices/internal/models"
	"github.com/google/uuid"
)

// JWTService defines the interface for JWT token operations
type JWTService interface {
	GenerateAccessToken(user *models.User) (string, error)
	GenerateRefreshToken(user *models.User) (string, error)
	ValidateToken(token string) (*JWTClaims, error)
	ExtractUserID(tokenString string) (uuid.UUID, error)
	ExtractTokenType(tokenString string) (string, error)
	GetTokenExpiry(tokenType string) time.Duration
}

// PasswordService defines the interface for password operations
type PasswordService interface {
	HashPassword(password string) (string, error)
	VerifyPassword(hashedPassword, password string) error
	IsStrongPassword(password string) (bool, string)
}

