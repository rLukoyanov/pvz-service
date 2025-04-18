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

// MockReceptionService is a mock implementation of the Reception service
type MockReceptionService struct {
	mock.Mock
}

func (m *MockReceptionService) GetActiveReceptionByPVZID(ctx context.Context, pvzID string) (*models.Reception, error) {
	args := m.Called(ctx, pvzID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Reception), args.Error(1)
}

func (m *MockReceptionService) CreateReception(ctx context.Context, reception models.Reception) error {
	args := m.Called(ctx, reception)
	return args.Error(0)
}

func setupReceptionEcho() (*echo.Echo, *MockReceptionService, *ReceptionHandler) {
	e := echo.New()
	mockService := new(MockReceptionService)
	s := &services.Services{
		ReceptionService: &services.ReceptionService{},
	}
	handler := NewReceptionHandler(s)
	return e, mockService, handler
}

func TestReceptionHandler_Create(t *testing.T) {
	e, mockService, handler := setupReceptionEcho()

	// Test case 1: Successful reception creation
	t.Run("successful reception creation", func(t *testing.T) {
		reqBody := map[string]string{
			"pvzId": "1",
		}
		reqJSON, _ := json.Marshal(reqBody)

		mockService.On("GetActiveReceptionByPVZID", mock.Anything, "1").
			Return(nil, nil)
		mockService.On("CreateReception", mock.Anything, mock.Anything).
			Return(nil)

		req := httptest.NewRequest(http.MethodPost, "/receptions", strings.NewReader(string(reqJSON)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.Create(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, rec.Code)

		var response map[string]string
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "created", response["message"])
	})

	// Test case 2: Invalid request body
	t.Run("invalid request body", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/receptions", strings.NewReader("invalid json"))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.Create(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	// Test case 3: Active reception exists
	t.Run("active reception exists", func(t *testing.T) {
		reqBody := map[string]string{
			"pvzId": "1",
		}
		reqJSON, _ := json.Marshal(reqBody)

		activeReception := &models.Reception{
			PvzId:  "1",
			Status: "in_progress",
		}
		mockService.On("GetActiveReceptionByPVZID", mock.Anything, "1").
			Return(activeReception, nil)

		req := httptest.NewRequest(http.MethodPost, "/receptions", strings.NewReader(string(reqJSON)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.Create(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	// Test case 4: Service error
	t.Run("service error", func(t *testing.T) {
		reqBody := map[string]string{
			"pvzId": "1",
		}
		reqJSON, _ := json.Marshal(reqBody)

		mockService.On("GetActiveReceptionByPVZID", mock.Anything, "1").
			Return(nil, assert.AnError)

		req := httptest.NewRequest(http.MethodPost, "/receptions", strings.NewReader(string(reqJSON)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.Create(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}
