package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/galafis/go-data-api-microservices/internal/config"
	"github.com/galafis/go-data-api-microservices/internal/models"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

// JWTClaims represents the claims in a JWT token
type JWTClaims struct {
	UserID    string      `json:"user_id"`
	Email     string      `json:"email"`
	Role      models.Role `json:"role"`
	TokenType string      `json:"token_type"`
	jwt.RegisteredClaims
}

// JWTService handles JWT token operations
type JWTService struct {
	config *config.AuthConfig
}

// NewJWTService creates a new JWT service
func NewJWTService(config *config.AuthConfig) *JWTService {
	return &JWTService{
		config: config,
	}
}

// GenerateAccessToken generates a new access token
func (s *JWTService) GenerateAccessToken(user *models.User) (string, error) {
	expirationTime := time.Now().Add(s.config.AccessTokenExpiry)
	
	claims := JWTClaims{
		UserID:    user.ID.String(),
		Email:     user.Email,
		Role:      user.Role,
		TokenType: "access",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "go-data-api",
			Subject:   user.ID.String(),
		},
	}
	
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.config.JWTSecret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}
	
	return tokenString, nil
}

// GenerateRefreshToken generates a new refresh token
func (s *JWTService) GenerateRefreshToken(user *models.User) (string, error) {
	expirationTime := time.Now().Add(s.config.RefreshTokenExpiry)
	
	claims := JWTClaims{
		UserID:    user.ID.String(),
		Email:     user.Email,
		Role:      user.Role,
		TokenType: "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "go-data-api",
			Subject:   user.ID.String(),
		},
	}
	
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.config.JWTSecret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}
	
	return tokenString, nil
}

// ValidateToken validates a JWT token
func (s *JWTService) ValidateToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.config.JWTSecret), nil
	})
	
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}
	
	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}
	
	return nil, errors.New("invalid token")
}

// ExtractUserID extracts the user ID from a JWT token
func (s *JWTService) ExtractUserID(tokenString string) (uuid.UUID, error) {
	claims, err := s.ValidateToken(tokenString)
	if err != nil {
		return uuid.Nil, err
	}
	
	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid user ID in token: %w", err)
	}
	
	return userID, nil
}

// ExtractTokenType extracts the token type from a JWT token
func (s *JWTService) ExtractTokenType(tokenString string) (string, error) {
	claims, err := s.ValidateToken(tokenString)
	if err != nil {
		return "", err
	}
	
	return claims.TokenType, nil
}

// GetTokenExpiry returns the expiry time of a token
func (s *JWTService) GetTokenExpiry(tokenType string) time.Duration {
	switch tokenType {
	case "access":
		return s.config.AccessTokenExpiry
	case "refresh":
		return s.config.RefreshTokenExpiry
	default:
		return 0
	}
}

