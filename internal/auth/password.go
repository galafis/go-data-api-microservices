package auth

import (
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// PasswordService handles password operations
type PasswordService struct {
	hashCost int
}

// NewPasswordService creates a new password service
func NewPasswordService(hashCost int) *PasswordService {
	if hashCost < bcrypt.MinCost {
		hashCost = bcrypt.DefaultCost
	}
	
	return &PasswordService{
		hashCost: hashCost,
	}
}

// HashPassword hashes a password
func (s *PasswordService) HashPassword(password string) (string, error) {
	if password == "" {
		return "", errors.New("password cannot be empty")
	}
	
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), s.hashCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	
	return string(hashedPassword), nil
}

// VerifyPassword verifies a password against a hash
func (s *PasswordService) VerifyPassword(hashedPassword, password string) error {
	if hashedPassword == "" || password == "" {
		return errors.New("password or hash cannot be empty")
	}
	
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// IsStrongPassword checks if a password is strong enough
func (s *PasswordService) IsStrongPassword(password string) (bool, string) {
	if len(password) < 8 {
		return false, "password must be at least 8 characters long"
	}
	
	var hasUpper, hasLower, hasNumber, hasSpecial bool
	for _, char := range password {
		switch {
		case 'A' <= char && char <= 'Z':
			hasUpper = true
		case 'a' <= char && char <= 'z':
			hasLower = true
		case '0' <= char && char <= '9':
			hasNumber = true
		case char == '!' || char == '@' || char == '#' || char == '$' || char == '%' || char == '^' || char == '&' || char == '*':
			hasSpecial = true
		}
	}
	
	if !hasUpper {
		return false, "password must contain at least one uppercase letter"
	}
	if !hasLower {
		return false, "password must contain at least one lowercase letter"
	}
	if !hasNumber {
		return false, "password must contain at least one number"
	}
	if !hasSpecial {
		return false, "password must contain at least one special character (!@#$%^&*)"
	}
	
	return true, ""
}

