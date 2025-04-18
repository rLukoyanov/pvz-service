package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"pvz-service/internal/models"
	"pvz-service/internal/services"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserService is a mock implementation of the User service
type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) CreateUser(ctx context.Context, user models.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserService) GetUserByEmail(ctx context.Context, email string) (models.User, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(models.User), args.Error(1)
}

func setupAuthEcho() (*echo.Echo, *MockUserService, *AuthHandler) {
	e := echo.New()
	mockService := new(MockUserService)
	s := &services.Services{
		UserService: &services.UserService{},
	}
	handler := NewAuthHandler(s)
	return e, mockService, handler
}

func TestAuthHandler_Register(t *testing.T) {
	e, mockService, handler := setupAuthEcho()

	// Test case 1: Successful registration
	t.Run("successful registration", func(t *testing.T) {
		reqBody := map[string]string{
			"email":    "test@example.com",
			"password": "password123",
			"role":     "client",
		}
		reqJSON, _ := json.Marshal(reqBody)

		mockService.On("CreateUser", mock.Anything, mock.Anything).
			Return(nil)

		req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(string(reqJSON)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.Register(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, rec.Code)

		var response map[string]string
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "test@example.com", response["email"])
		assert.Equal(t, "client", response["role"])
	})

	// Test case 2: Invalid request body
	t.Run("invalid request body", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader("invalid json"))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.Register(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	// Test case 3: Invalid role
	t.Run("invalid role", func(t *testing.T) {
		reqBody := map[string]string{
			"email":    "test@example.com",
			"password": "password123",
			"role":     "invalid_role",
		}
		reqJSON, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(string(reqJSON)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.Register(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	// Test case 4: Service error
	t.Run("service error", func(t *testing.T) {
		reqBody := map[string]string{
			"email":    "test@example.com",
			"password": "password123",
			"role":     "client",
		}
		reqJSON, _ := json.Marshal(reqBody)

		mockService.On("CreateUser", mock.Anything, mock.Anything).
			Return(assert.AnError)

		req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(string(reqJSON)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.Register(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

func TestAuthHandler_Login(t *testing.T) {
	e, mockService, handler := setupAuthEcho()

	// Test case 1: Successful login
	t.Run("successful login", func(t *testing.T) {
		reqBody := map[string]string{
			"email":    "test@example.com",
			"password": "password123",
		}
		reqJSON, _ := json.Marshal(reqBody)

		user := models.User{
			Email:    "test@example.com",
			Password: "$2a$10$hashedpassword", // bcrypt hash of "password123"
			Role:     "client",
		}
		mockService.On("GetUserByEmail", mock.Anything, "test@example.com").
			Return(user, nil)

		req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(string(reqJSON)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.Login(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var response map[string]string
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotEmpty(t, response["token"])
	})

	// Test case 2: Invalid request body
	t.Run("invalid request body", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader("invalid json"))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.Login(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	// Test case 3: User not found
	t.Run("user not found", func(t *testing.T) {
		reqBody := map[string]string{
			"email":    "nonexistent@example.com",
			"password": "password123",
		}
		reqJSON, _ := json.Marshal(reqBody)

		mockService.On("GetUserByEmail", mock.Anything, "nonexistent@example.com").
			Return(models.User{}, assert.AnError)

		req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(string(reqJSON)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.Login(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	// Test case 4: Invalid password
	t.Run("invalid password", func(t *testing.T) {
		reqBody := map[string]string{
			"email":    "test@example.com",
			"password": "wrongpassword",
		}
		reqJSON, _ := json.Marshal(reqBody)

		user := models.User{
			Email:    "test@example.com",
			Password: "$2a$10$hashedpassword", // bcrypt hash of "password123"
			Role:     "client",
		}
		mockService.On("GetUserByEmail", mock.Anything, "test@example.com").
			Return(user, nil)

		req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(string(reqJSON)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.Login(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})
}
