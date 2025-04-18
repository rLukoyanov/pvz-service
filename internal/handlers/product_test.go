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

// MockProductService is a mock implementation of the Product service
type MockProductService struct {
	mock.Mock
}

func (m *MockProductService) AddProduct(ctx context.Context, product models.Product, pvzId string) error {
	args := m.Called(ctx, product, pvzId)
	return args.Error(0)
}

func (m *MockProductService) DeleteLastProduct(ctx context.Context, pvzID string) error {
	args := m.Called(ctx, pvzID)
	return args.Error(0)
}

func setupProductEcho() (*echo.Echo, *MockProductService, *ItemHandler) {
	e := echo.New()
	mockService := new(MockProductService)
	s := &services.Services{
		ProductService: &services.ProductService{},
	}
	handler := NewProductHandler(s)
	return e, mockService, handler
}

func TestItemHandler_AddProduct(t *testing.T) {
	e, mockService, handler := setupProductEcho()

	// Test case 1: Successful product addition
	t.Run("successful product addition", func(t *testing.T) {
		reqBody := map[string]string{
			"type":  "electronics",
			"PvzId": "1",
		}
		reqJSON, _ := json.Marshal(reqBody)

		mockService.On("AddProduct", mock.Anything, mock.Anything, "1").
			Return(nil)

		req := httptest.NewRequest(http.MethodPost, "/products", strings.NewReader(string(reqJSON)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.AddProduct(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, rec.Code)

		var response map[string]string
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "created", response["message"])
	})

	// Test case 2: Invalid request body
	t.Run("invalid request body", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/products", strings.NewReader("invalid json"))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.AddProduct(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	// Test case 3: Missing required fields
	t.Run("missing required fields", func(t *testing.T) {
		reqBody := map[string]string{
			"type": "electronics",
		}
		reqJSON, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/products", strings.NewReader(string(reqJSON)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.AddProduct(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	// Test case 4: Service error
	t.Run("service error", func(t *testing.T) {
		reqBody := map[string]string{
			"type":  "electronics",
			"PvzId": "1",
		}
		reqJSON, _ := json.Marshal(reqBody)

		mockService.On("AddProduct", mock.Anything, mock.Anything, "1").
			Return(assert.AnError)

		req := httptest.NewRequest(http.MethodPost, "/products", strings.NewReader(string(reqJSON)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.AddProduct(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}
