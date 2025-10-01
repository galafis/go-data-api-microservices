package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/galafis/go-data-api-microservices/internal/auth"
	"github.com/galafis/go-data-api-microservices/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockJWTService é um mock para auth.JWTService
type MockJWTService struct {
	mock.Mock
}

func (m *MockJWTService) GenerateAccessToken(user *models.User) (string, error) {
	args := m.Called(user)
	return args.String(0), args.Error(1)
}

func (m *MockJWTService) GenerateRefreshToken(user *models.User) (string, error) {
	args := m.Called(user)
	return args.String(0), args.Error(1)
}

func (m *MockJWTService) ValidateToken(token string) (*auth.JWTClaims, error) {
	args := m.Called(token)
	return args.Get(0).(*auth.JWTClaims), args.Error(1)
}

func (m *MockJWTService) ExtractUserID(tokenString string) (uuid.UUID, error) {
	args := m.Called(tokenString)
	return args.Get(0).(uuid.UUID), args.Error(1)
}

func (m *MockJWTService) ExtractTokenType(tokenString string) (string, error) {
	args := m.Called(tokenString)
	return args.String(0), args.Error(1)
}

func (m *MockJWTService) GetTokenExpiry(tokenType string) time.Duration {
	args := m.Called(tokenType)
	return args.Get(0).(time.Duration)
}

// MockPasswordService é um mock para auth.PasswordService
type MockPasswordService struct {
	mock.Mock
}

func (m *MockPasswordService) HashPassword(password string) (string, error) {
	args := m.Called(password)
	return args.String(0), args.Error(1)
}

func (m *MockPasswordService) VerifyPassword(hashedPassword, password string) error {
	args := m.Called(hashedPassword, password)
	return args.Error(0)
}

func (m *MockPasswordService) IsStrongPassword(password string) (bool, string) {
	args := m.Called(password)
	return args.Bool(0), args.String(1)
}

// MockUserRepository é um mock para UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) FindByEmail(email string) (*models.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) FindByID(userID uuid.UUID) (*models.User, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) Create(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) Update(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(userID uuid.UUID) error {
	args := m.Called(userID)
	return args.Error(0)
}

func (m *MockUserRepository) UpdateRefreshToken(userID uuid.UUID, token string) error {
	args := m.Called(userID, token)
	return args.Error(0)
}

func (m *MockUserRepository) ClearRefreshToken(userID uuid.UUID) error {
	args := m.Called(userID)
	return args.Error(0)
}

func TestAuthHandler_Register(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockJWTService := new(MockJWTService)
	mockPasswordService := new(MockPasswordService)
	mockUserRepository := new(MockUserRepository)

	authHandler := NewAuthHandler(mockJWTService, mockPasswordService, mockUserRepository)

	_ = authHandler // Para evitar erro de variável não utilizada

	r := gin.Default()
	r.POST("/register", authHandler.Register)

	// Test Case 1: Successful Registration
	t.Run("Successful Registration", func(t *testing.T) {
		registerReq := models.RegisterRequest{
			Email:     "test@example.com",
			Password:  "StrongPassword123!",
			FirstName: "John",
			LastName:  "Doe",
		}
		jsonValue, _ := json.Marshal(registerReq)

		mockUserRepository.On("FindByEmail", registerReq.Email).Return(nil, nil).Once()
		mockPasswordService.On("IsStrongPassword", registerReq.Password).Return(true, "").Once()
		mockPasswordService.On("HashPassword", registerReq.Password).Return("hashedpassword", nil).Once()
		mockUserRepository.On("Create", mock.AnythingOfType("*models.User")).Return(nil).Once()
		mockJWTService.On("GenerateAccessToken", mock.AnythingOfType("*models.User")).Return("accesstoken", nil).Once()
		mockJWTService.On("GenerateRefreshToken", mock.AnythingOfType("*models.User")).Return("refreshtoken", nil).Once()
		mockUserRepository.On("UpdateRefreshToken", mock.AnythingOfType("uuid.UUID"), "refreshtoken").Return(nil).Once()
		mockJWTService.On("GetTokenExpiry", "access").Return(15 * time.Minute).Once()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		var res models.AuthResponse
		json.Unmarshal(w.Body.Bytes(), &res)
		assert.Equal(t, "test@example.com", res.User.Email)
		assert.Equal(t, "accesstoken", res.AccessToken)
		assert.Equal(t, "refreshtoken", res.RefreshToken)

		mockUserRepository.AssertExpectations(t)
		mockPasswordService.AssertExpectations(t)
		mockJWTService.AssertExpectations(t)
	})

	// Test Case 2: Registration with existing email
	t.Run("Existing Email", func(t *testing.T) {
		registerReq := models.RegisterRequest{
			Email:     "existing@example.com",
			Password:  "StrongPassword123!",
			FirstName: "Jane",
			LastName:  "Doe",
		}
		jsonValue, _ := json.Marshal(registerReq)

		existingUser := &models.User{ID: uuid.New(), Email: "existing@example.com"}
		mockUserRepository.On("FindByEmail", registerReq.Email).Return(existingUser, nil).Once()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusConflict, w.Code)
		assert.Contains(t, w.Body.String(), "Email already exists")

		mockUserRepository.AssertExpectations(t)
		mockPasswordService.AssertNotCalled(t, "IsStrongPassword")
		mockJWTService.AssertNotCalled(t, "GenerateAccessToken")
	})

	// Test Case 3: Registration with weak password
	t.Run("Weak Password", func(t *testing.T) {
		registerReq := models.RegisterRequest{
			Email:     "weakpass@example.com",
			Password:  "weak",
			FirstName: "Peter",
			LastName:  "Pan",
		}
		jsonValue, _ := json.Marshal(registerReq)

		mockUserRepository.On("FindByEmail", registerReq.Email).Return(nil, nil).Once()
		mockPasswordService.On("IsStrongPassword", registerReq.Password).Return(false, "Password is too weak").Once()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Password is too weak")

		mockUserRepository.AssertExpectations(t)
		mockPasswordService.AssertExpectations(t)
		mockJWTService.AssertNotCalled(t, "GenerateAccessToken")
	})

	// Test Case 4: Error hashing password
	t.Run("Error Hashing Password", func(t *testing.T) {
		registerReq := models.RegisterRequest{
			Email:     "hashfail@example.com",
			Password:  "StrongPassword123!",
			FirstName: "Alice",
			LastName:  "Wonderland",
		}
		jsonValue, _ := json.Marshal(registerReq)

		mockUserRepository.On("FindByEmail", registerReq.Email).Return(nil, nil).Once()
		mockPasswordService.On("IsStrongPassword", registerReq.Password).Return(true, "").Once()
		mockPasswordService.On("HashPassword", registerReq.Password).Return("", errors.New("hashing error")).Once()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "Internal server error")

		mockUserRepository.AssertExpectations(t)
		mockPasswordService.AssertExpectations(t)
		mockJWTService.AssertNotCalled(t, "GenerateAccessToken")
	})

	// Test Case 5: Error creating user
	t.Run("Error Creating User", func(t *testing.T) {
		registerReq := models.RegisterRequest{
			Email:     "createfail@example.com",
			Password:  "StrongPassword123!",
			FirstName: "Bob",
			LastName:  "Builder",
		}
		jsonValue, _ := json.Marshal(registerReq)

		mockUserRepository.On("FindByEmail", registerReq.Email).Return(nil, nil).Once()
		mockPasswordService.On("IsStrongPassword", registerReq.Password).Return(true, "").Once()
		mockPasswordService.On("HashPassword", registerReq.Password).Return("hashedpassword", nil).Once()
		mockUserRepository.On("Create", mock.AnythingOfType("*models.User")).Return(errors.New("create error")).Once()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "Internal server error")

		mockUserRepository.AssertExpectations(t)
		mockPasswordService.AssertExpectations(t)
		mockJWTService.AssertNotCalled(t, "GenerateAccessToken")
	})

	// Test Case 6: Error generating access token
	t.Run("Error Generating Access Token", func(t *testing.T) {
		registerReq := models.RegisterRequest{
			Email:     "accesstokenfail@example.com",
			Password:  "StrongPassword123!",
			FirstName: "Charlie",
			LastName:  "Chaplin",
		}
		jsonValue, _ := json.Marshal(registerReq)

		mockUserRepository.On("FindByEmail", registerReq.Email).Return(nil, nil).Once()
		mockPasswordService.On("IsStrongPassword", registerReq.Password).Return(true, "").Once()
		mockPasswordService.On("HashPassword", registerReq.Password).Return("hashedpassword", nil).Once()
		mockUserRepository.On("Create", mock.AnythingOfType("*models.User")).Return(nil).Once()
		mockJWTService.On("GenerateAccessToken", mock.AnythingOfType("*models.User")).Return("", errors.New("access token error")).Once()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "Internal server error")

		mockUserRepository.AssertExpectations(t)
		mockPasswordService.AssertExpectations(t)
		mockJWTService.AssertExpectations(t)
	})

	// Test Case 7: Error generating refresh token
	t.Run("Error Generating Refresh Token", func(t *testing.T) {
		registerReq := models.RegisterRequest{
			Email:     "refreshtokenfail@example.com",
			Password:  "StrongPassword123!",
			FirstName: "David",
			LastName:  "Copperfield",
		}
		jsonValue, _ := json.Marshal(registerReq)

		mockUserRepository.On("FindByEmail", registerReq.Email).Return(nil, nil).Once()
		mockPasswordService.On("IsStrongPassword", registerReq.Password).Return(true, "").Once()
		mockPasswordService.On("HashPassword", registerReq.Password).Return("hashedpassword", nil).Once()
		mockUserRepository.On("Create", mock.AnythingOfType("*models.User")).Return(nil).Once()
		mockJWTService.On("GenerateAccessToken", mock.AnythingOfType("*models.User")).Return("accesstoken", nil).Once()
		mockJWTService.On("GenerateRefreshToken", mock.AnythingOfType("*models.User")).Return("", errors.New("refresh token error")).Once()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "Internal server error")

		mockUserRepository.AssertExpectations(t)
		mockPasswordService.AssertExpectations(t)
		mockJWTService.AssertExpectations(t)
	})

	// Test Case 8: Error updating refresh token
	t.Run("Error Updating Refresh Token", func(t *testing.T) {
		registerReq := models.RegisterRequest{
			Email:     "updaterefreshfail@example.com",
			Password:  "StrongPassword123!",
			FirstName: "Eve",
			LastName:  "Adams",
		}
		jsonValue, _ := json.Marshal(registerReq)

		mockUserRepository.On("FindByEmail", registerReq.Email).Return(nil, nil).Once()
		mockPasswordService.On("IsStrongPassword", registerReq.Password).Return(true, "").Once()
		mockPasswordService.On("HashPassword", registerReq.Password).Return("hashedpassword", nil).Once()
		mockUserRepository.On("Create", mock.AnythingOfType("*models.User")).Return(nil).Once()
		mockJWTService.On("GenerateAccessToken", mock.AnythingOfType("*models.User")).Return("accesstoken", nil).Once()
		mockJWTService.On("GenerateRefreshToken", mock.AnythingOfType("*models.User")).Return("refreshtoken", nil).Once()
		mockUserRepository.On("UpdateRefreshToken", mock.AnythingOfType("uuid.UUID"), "refreshtoken").Return(errors.New("update refresh token error")).Once()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "Internal server error")

		mockUserRepository.AssertExpectations(t)
		mockPasswordService.AssertExpectations(t)
		mockJWTService.AssertExpectations(t)
	})
}

func TestAuthHandler_Login(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockJWTService := new(MockJWTService)
	mockPasswordService := new(MockPasswordService)
	mockUserRepository := new(MockUserRepository)

	authHandler := NewAuthHandler(mockJWTService, mockPasswordService, mockUserRepository)

	_ = authHandler // Para evitar erro de variável não utilizada

	r := gin.Default()
	r.POST("/login", authHandler.Login)

	// Test Case 1: Successful Login
	t.Run("Successful Login", func(t *testing.T) {
		loginReq := models.LoginRequest{
			Email:    "test@example.com",
			Password: "password123",
		}
		jsonValue, _ := json.Marshal(loginReq)

		user := &models.User{
			ID:        uuid.New(),
			Email:     "test@example.com",
			Password:  "hashedpassword",
			FirstName: "John",
			LastName:  "Doe",
			Active:    true,
		}

		mockUserRepository.On("FindByEmail", loginReq.Email).Return(user, nil).Once()
		mockPasswordService.On("VerifyPassword", user.Password, loginReq.Password).Return(nil).Once()
		mockJWTService.On("GenerateAccessToken", user).Return("accesstoken", nil).Once()
		mockJWTService.On("GenerateRefreshToken", user).Return("refreshtoken", nil).Once()
		mockUserRepository.On("UpdateRefreshToken", user.ID, "refreshtoken").Return(nil).Once()
		mockJWTService.On("GetTokenExpiry", "access").Return(15 * time.Minute).Once()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var res models.AuthResponse
		json.Unmarshal(w.Body.Bytes(), &res)
		assert.Equal(t, "test@example.com", res.User.Email)
		assert.Equal(t, "accesstoken", res.AccessToken)
		assert.Equal(t, "refreshtoken", res.RefreshToken)

		mockUserRepository.AssertExpectations(t)
		mockPasswordService.AssertExpectations(t)
		mockJWTService.AssertExpectations(t)
	})

	// Test Case 2: Login with invalid credentials (user not found)
	t.Run("Invalid Credentials - User Not Found", func(t *testing.T) {
		loginReq := models.LoginRequest{
			Email:    "nonexistent@example.com",
			Password: "password123",
		}
		jsonValue, _ := json.Marshal(loginReq)

		mockUserRepository.On("FindByEmail", loginReq.Email).Return(nil, nil).Once()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid email or password")

		mockUserRepository.AssertExpectations(t)
		mockPasswordService.AssertNotCalled(t, "VerifyPassword")
		mockJWTService.AssertNotCalled(t, "GenerateAccessToken")
	})

	// Test Case 3: Login with inactive account
	t.Run("Inactive Account", func(t *testing.T) {
		loginReq := models.LoginRequest{
			Email:    "inactive@example.com",
			Password: "password123",
		}
		jsonValue, _ := json.Marshal(loginReq)

		user := &models.User{
			ID:        uuid.New(),
			Email:     "inactive@example.com",
			Password:  "hashedpassword",
			FirstName: "Inactive",
			LastName:  "User",
			Active:    false,
		}

		mockUserRepository.On("FindByEmail", loginReq.Email).Return(user, nil).Once()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "Account is inactive")

		mockUserRepository.AssertExpectations(t)
		mockPasswordService.AssertNotCalled(t, "VerifyPassword")
		mockJWTService.AssertNotCalled(t, "GenerateAccessToken")
	})

	// Test Case 4: Login with invalid password
	t.Run("Invalid Password", func(t *testing.T) {
		loginReq := models.LoginRequest{
			Email:    "test@example.com",
			Password: "wrongpassword",
		}
		jsonValue, _ := json.Marshal(loginReq)

		user := &models.User{
			ID:        uuid.New(),
			Email:     "test@example.com",
			Password:  "hashedpassword",
			FirstName: "John",
			LastName:  "Doe",
			Active:    true,
		}

		mockUserRepository.On("FindByEmail", loginReq.Email).Return(user, nil).Once()
		mockPasswordService.On("VerifyPassword", user.Password, loginReq.Password).Return(errors.New("invalid password")).Once()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid email or password")

		mockUserRepository.AssertExpectations(t)
		mockPasswordService.AssertExpectations(t)
		mockJWTService.AssertNotCalled(t, "GenerateAccessToken")
	})

	// Test Case 5: Error generating access token during login
	t.Run("Error Generating Access Token on Login", func(t *testing.T) {
		loginReq := models.LoginRequest{
			Email:    "test@example.com",
			Password: "password123",
		}
		jsonValue, _ := json.Marshal(loginReq)

		user := &models.User{
			ID:        uuid.New(),
			Email:     "test@example.com",
			Password:  "hashedpassword",
			FirstName: "John",
			LastName:  "Doe",
			Active:    true,
		}

		mockUserRepository.On("FindByEmail", loginReq.Email).Return(user, nil).Once()
		mockPasswordService.On("VerifyPassword", user.Password, loginReq.Password).Return(nil).Once()
		mockJWTService.On("GenerateAccessToken", user).Return("", errors.New("access token error")).Once()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "Internal server error")

		mockUserRepository.AssertExpectations(t)
		mockPasswordService.AssertExpectations(t)
		mockJWTService.AssertExpectations(t)
	})

	// Test Case 6: Error generating refresh token during login
	t.Run("Error Generating Refresh Token on Login", func(t *testing.T) {
		loginReq := models.LoginRequest{
			Email:    "test@example.com",
			Password: "password123",
		}
		jsonValue, _ := json.Marshal(loginReq)

		user := &models.User{
			ID:        uuid.New(),
			Email:     "test@example.com",
			Password:  "hashedpassword",
			FirstName: "John",
			LastName:  "Doe",
			Active:    true,
		}

		mockUserRepository.On("FindByEmail", loginReq.Email).Return(user, nil).Once()
		mockPasswordService.On("VerifyPassword", user.Password, loginReq.Password).Return(nil).Once()
		mockJWTService.On("GenerateAccessToken", user).Return("accesstoken", nil).Once()
		mockJWTService.On("GenerateRefreshToken", user).Return("", errors.New("refresh token error")).Once()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "Internal server error")

		mockUserRepository.AssertExpectations(t)
		mockPasswordService.AssertExpectations(t)
		mockJWTService.AssertExpectations(t)
	})

	// Test Case 7: Error updating refresh token during login
	t.Run("Error Updating Refresh Token on Login", func(t *testing.T) {
		loginReq := models.LoginRequest{
			Email:    "test@example.com",
			Password: "password123",
		}
		jsonValue, _ := json.Marshal(loginReq)

		user := &models.User{
			ID:        uuid.New(),
			Email:     "test@example.com",
			Password:  "hashedpassword",
			FirstName: "John",
			LastName:  "Doe",
			Active:    true,
		}

		mockUserRepository.On("FindByEmail", loginReq.Email).Return(user, nil).Once()
		mockPasswordService.On("VerifyPassword", user.Password, loginReq.Password).Return(nil).Once()
		mockJWTService.On("GenerateAccessToken", user).Return("accesstoken", nil).Once()
		mockJWTService.On("GenerateRefreshToken", user).Return("refreshtoken", nil).Once()
		mockUserRepository.On("UpdateRefreshToken", user.ID, "refreshtoken").Return(errors.New("update refresh token error")).Once()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "Internal server error")

		mockUserRepository.AssertExpectations(t)
		mockPasswordService.AssertExpectations(t)
		mockJWTService.AssertExpectations(t)
	})
}

func TestAuthHandler_RefreshToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockJWTService := new(MockJWTService)
	mockPasswordService := new(MockPasswordService)
	mockUserRepository := new(MockUserRepository)

	authHandler := NewAuthHandler(mockJWTService, mockPasswordService, mockUserRepository)

	_ = authHandler // Para evitar erro de variável não utilizada

	r := gin.Default()
	r.POST("/refresh", authHandler.RefreshToken)

	// Test Case 1: Successful Token Refresh
	t.Run("Successful Token Refresh", func(t *testing.T) {
		refreshTokenReq := models.RefreshTokenRequest{
			RefreshToken: "validrefreshtoken",
		}
		jsonValue, _ := json.Marshal(refreshTokenReq)

		userID := uuid.New()
		claims := &auth.JWTClaims{UserID: userID.String(), TokenType: "refresh"}
		user := &models.User{
			ID:           userID,
			Email:        "test@example.com",
			FirstName:    "John",
			LastName:     "Doe",
			Active:       true,
			RefreshToken: "validrefreshtoken",
		}

		mockJWTService.On("ValidateToken", refreshTokenReq.RefreshToken).Return(claims, nil).Once()
		mockUserRepository.On("FindByID", userID).Return(user, nil).Once()
		mockJWTService.On("GenerateAccessToken", user).Return("newaccesstoken", nil).Once()
		mockJWTService.On("GenerateRefreshToken", user).Return("newrefreshtoken", nil).Once()
		mockUserRepository.On("UpdateRefreshToken", userID, "newrefreshtoken").Return(nil).Once()
		mockJWTService.On("GetTokenExpiry", "access").Return(15 * time.Minute).Once()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/refresh", bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var res models.AuthResponse
		json.Unmarshal(w.Body.Bytes(), &res)
		assert.Equal(t, "test@example.com", res.User.Email)
		assert.Equal(t, "newaccesstoken", res.AccessToken)
		assert.Equal(t, "newrefreshtoken", res.RefreshToken)

		mockJWTService.AssertExpectations(t)
		mockUserRepository.AssertExpectations(t)
	})

	// Test Case 2: Invalid Refresh Token
	t.Run("Invalid Refresh Token", func(t *testing.T) {
		refreshTokenReq := models.RefreshTokenRequest{
			RefreshToken: "invalidtoken",
		}
		jsonValue, _ := json.Marshal(refreshTokenReq)

		mockJWTService.On("ValidateToken", refreshTokenReq.RefreshToken).Return(&auth.JWTClaims{}, errors.New("invalid token")).Once()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/refresh", bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid refresh token")

		mockJWTService.AssertExpectations(t)
		mockUserRepository.AssertNotCalled(t, "FindByID")
	})

	// Test Case 3: Incorrect Token Type
	t.Run("Incorrect Token Type", func(t *testing.T) {
		refreshTokenReq := models.RefreshTokenRequest{
			RefreshToken: "accesstoken",
		}
		jsonValue, _ := json.Marshal(refreshTokenReq)

		claims := &auth.JWTClaims{UserID: uuid.New().String(), TokenType: "access"}
		mockJWTService.On("ValidateToken", refreshTokenReq.RefreshToken).Return(claims, nil).Once()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/refresh", bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid token type")

		mockJWTService.AssertExpectations(t)
		mockUserRepository.AssertNotCalled(t, "FindByID")
	})

	// Test Case 4: User Not Found during refresh
	t.Run("User Not Found", func(t *testing.T) {
		refreshTokenReq := models.RefreshTokenRequest{
			RefreshToken: "validrefreshtoken",
		}
		jsonValue, _ := json.Marshal(refreshTokenReq)

		userID := uuid.New()
		claims := &auth.JWTClaims{UserID: userID.String(), TokenType: "refresh"}
		mockJWTService.On("ValidateToken", refreshTokenReq.RefreshToken).Return(claims, nil).Once()
		mockUserRepository.On("FindByID", userID).Return(nil, nil).Once()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/refresh", bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "User not found")

		mockJWTService.AssertExpectations(t)
		mockUserRepository.AssertExpectations(t)
	})

	// Test Case 5: Inactive User during refresh
	t.Run("Inactive User", func(t *testing.T) {
		refreshTokenReq := models.RefreshTokenRequest{
			RefreshToken: "validrefreshtoken",
		}
		jsonValue, _ := json.Marshal(refreshTokenReq)

		userID := uuid.New()
		claims := &auth.JWTClaims{UserID: userID.String(), TokenType: "refresh"}
		user := &models.User{
			ID:           userID,
			Email:        "test@example.com",
			Active:       false,
			RefreshToken: "validrefreshtoken",
		}

		mockJWTService.On("ValidateToken", refreshTokenReq.RefreshToken).Return(claims, nil).Once()
		mockUserRepository.On("FindByID", userID).Return(user, nil).Once()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/refresh", bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "Account is inactive")

		mockJWTService.AssertExpectations(t)
		mockUserRepository.AssertExpectations(t)
	})

	// Test Case 6: Mismatched Refresh Token
	t.Run("Mismatched Refresh Token", func(t *testing.T) {
		refreshTokenReq := models.RefreshTokenRequest{
			RefreshToken: "validrefreshtoken",
		}
		jsonValue, _ := json.Marshal(refreshTokenReq)

		userID := uuid.New()
		claims := &auth.JWTClaims{UserID: userID.String(), TokenType: "refresh"}
		user := &models.User{
			ID:           userID,
			Email:        "test@example.com",
			Active:       true,
			RefreshToken: "anotherrefreshtoken", // Mismatched
		}

		mockJWTService.On("ValidateToken", refreshTokenReq.RefreshToken).Return(claims, nil).Once()
		mockUserRepository.On("FindByID", userID).Return(user, nil).Once()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/refresh", bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid refresh token")

		mockJWTService.AssertExpectations(t)
		mockUserRepository.AssertExpectations(t)
	})

	// Test Case 7: Error generating new access token
	t.Run("Error Generating New Access Token", func(t *testing.T) {
		refreshTokenReq := models.RefreshTokenRequest{
			RefreshToken: "validrefreshtoken",
		}
		jsonValue, _ := json.Marshal(refreshTokenReq)

		userID := uuid.New()
		claims := &auth.JWTClaims{UserID: userID.String(), TokenType: "refresh"}
		user := &models.User{
			ID:           userID,
			Email:        "test@example.com",
			FirstName:    "John",
			LastName:     "Doe",
			Active:       true,
			RefreshToken: "validrefreshtoken",
		}

		mockJWTService.On("ValidateToken", refreshTokenReq.RefreshToken).Return(claims, nil).Once()
		mockUserRepository.On("FindByID", userID).Return(user, nil).Once()
		mockJWTService.On("GenerateAccessToken", user).Return("", errors.New("access token error")).Once()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/refresh", bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "Internal server error")

		mockJWTService.AssertExpectations(t)
		mockUserRepository.AssertExpectations(t)
	})

	// Test Case 8: Error generating new refresh token
	t.Run("Error Generating New Refresh Token", func(t *testing.T) {
		refreshTokenReq := models.RefreshTokenRequest{
			RefreshToken: "validrefreshtoken",
		}
		jsonValue, _ := json.Marshal(refreshTokenReq)

		userID := uuid.New()
		claims := &auth.JWTClaims{UserID: userID.String(), TokenType: "refresh"}
		user := &models.User{
			ID:           userID,
			Email:        "test@example.com",
			FirstName:    "John",
			LastName:     "Doe",
			Active:       true,
			RefreshToken: "validrefreshtoken",
		}

		mockJWTService.On("ValidateToken", refreshTokenReq.RefreshToken).Return(claims, nil).Once()
		mockUserRepository.On("FindByID", userID).Return(user, nil).Once()
		mockJWTService.On("GenerateAccessToken", user).Return("newaccesstoken", nil).Once()
		mockJWTService.On("GenerateRefreshToken", user).Return("", errors.New("refresh token error")).Once()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/refresh", bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "Internal server error")

		mockJWTService.AssertExpectations(t)
		mockUserRepository.AssertExpectations(t)
	})

	// Test Case 9: Error updating refresh token
	t.Run("Error Updating Refresh Token", func(t *testing.T) {
		refreshTokenReq := models.RefreshTokenRequest{
			RefreshToken: "validrefreshtoken",
		}
		jsonValue, _ := json.Marshal(refreshTokenReq)

		userID := uuid.New()
		claims := &auth.JWTClaims{UserID: userID.String(), TokenType: "refresh"}
		user := &models.User{
			ID:           userID,
			Email:        "test@example.com",
			FirstName:    "John",
			LastName:     "Doe",
			Active:       true,
			RefreshToken: "validrefreshtoken",
		}

		mockJWTService.On("ValidateToken", refreshTokenReq.RefreshToken).Return(claims, nil).Once()
		mockUserRepository.On("FindByID", userID).Return(user, nil).Once()
		mockJWTService.On("GenerateAccessToken", user).Return("newaccesstoken", nil).Once()
		mockJWTService.On("GenerateRefreshToken", user).Return("newrefreshtoken", nil).Once()
		mockUserRepository.On("UpdateRefreshToken", userID, "newrefreshtoken").Return(errors.New("update refresh token error")).Once()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/refresh", bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "Internal server error")

		mockJWTService.AssertExpectations(t)
		mockUserRepository.AssertExpectations(t)
	})
}

func TestAuthHandler_Logout(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockJWTService := new(MockJWTService)
	mockPasswordService := new(MockPasswordService)
	mockUserRepository := new(MockUserRepository)

	authHandler := NewAuthHandler(mockJWTService, mockPasswordService, mockUserRepository)

	// Função auxiliar para configurar e executar o teste
	setupLogoutTest := func(userID *uuid.UUID) *httptest.ResponseRecorder {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/logout", nil)
		req.Header.Set("Content-Type", "application/json")

		router := gin.Default()
		router.POST("/logout", func(c *gin.Context) {
			if userID != nil {
				c.Set("user_id", *userID)
			}
			authHandler.Logout(c)
		})
		router.ServeHTTP(w, req)
		return w
	}

	// Test Case 1: Successful Logout
	t.Run("Successful Logout", func(t *testing.T) {
		userID := uuid.New()
		mockUserRepository.On("ClearRefreshToken", userID).Return(nil).Once()

		w := setupLogoutTest(&userID)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "Logged out successfully")

		mockUserRepository.AssertExpectations(t)
	})

	// Test Case 2: Unauthorized (user_id not in context)
	t.Run("Unauthorized - No User ID", func(t *testing.T) {
		w := setupLogoutTest(nil)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "Unauthorized")

		mockUserRepository.AssertNotCalled(t, "ClearRefreshToken")
	})

	// Test Case 3: Error clearing refresh token
	t.Run("Error Clearing Refresh Token", func(t *testing.T) {
		userID := uuid.New()
		mockUserRepository.On("ClearRefreshToken", userID).Return(errors.New("clear refresh token error")).Once()

		w := setupLogoutTest(&userID)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "Internal server error")

		mockUserRepository.AssertExpectations(t)
	})
}

