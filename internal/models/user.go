package models

import (
	"time"

	"github.com/google/uuid"
)

// Role represents a user role in the system
type Role string

const (
	RoleAdmin  Role = "admin"
	RoleUser   Role = "user"
	RoleViewer Role = "viewer"
)

// User represents a user in the system
type User struct {
	ID           uuid.UUID       `json:"id" bson:"_id"`
	Email        string          `json:"email" bson:"email"`
	Password     string          `json:"-" bson:"password"`
	FirstName    string          `json:"first_name" bson:"first_name"`
	LastName     string          `json:"last_name" bson:"last_name"`
	Role         Role            `json:"role" bson:"role"`
	Active       bool            `json:"active" bson:"active"`
	Verified     bool            `json:"verified" bson:"verified"`
	RefreshToken string          `json:"-" bson:"refresh_token,omitempty"`
	Metadata     map[string]any  `json:"metadata,omitempty" bson:"metadata,omitempty"`
	CreatedAt    time.Time       `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at" bson:"updated_at"`
	LastLoginAt  *time.Time      `json:"last_login_at,omitempty" bson:"last_login_at,omitempty"`
}

// RegisterRequest represents a user registration request
type RegisterRequest struct {
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=8"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
}

// LoginRequest represents a user login request
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// RefreshTokenRequest represents a token refresh request
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// UpdateUserRequest represents a request to update user information
type UpdateUserRequest struct {
	FirstName string         `json:"first_name,omitempty"`
	LastName  string         `json:"last_name,omitempty"`
	Password  string         `json:"password,omitempty" binding:"omitempty,min=8"`
	Metadata  map[string]any `json:"metadata,omitempty"`
}

// UserResponse represents a user response
type UserResponse struct {
	ID        uuid.UUID      `json:"id"`
	Email     string         `json:"email"`
	FirstName string         `json:"first_name"`
	LastName  string         `json:"last_name"`
	Role      Role           `json:"role"`
	Active    bool           `json:"active"`
	Verified  bool           `json:"verified"`
	Metadata  map[string]any `json:"metadata,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

// AuthResponse represents an authentication response
type AuthResponse struct {
	User        UserResponse `json:"user"`
	AccessToken string       `json:"access_token"`
	TokenType   string       `json:"token_type"`
	ExpiresIn   int64        `json:"expires_in"`
	RefreshToken string      `json:"refresh_token"`
}

