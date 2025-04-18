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
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockPVZService struct {
	mock.Mock
}

func (m *MockPVZService) GetAll(ctx context.Context, page, limit, from, to string) ([]models.FullPVZ, error) {
	args := m.Called(ctx, page, limit, from, to)
	return args.Get(0).([]models.FullPVZ), args.Error(1)
}

func (m *MockPVZService) CreatePVZ(ctx context.Context, pvz models.PVZ) (models.PVZ, error) {
	args := m.Called(ctx, pvz)
	return args.Get(0).(models.PVZ), args.Error(1)
}

func (m *MockPVZService) GetPVZByID(ctx context.Context, id string) (models.PVZ, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(models.PVZ), args.Error(1)
}

func (m *MockPVZService) DeletePVZ(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockPVZService) DeleteLastProduct(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockPVZService) CloseLastReception(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func setupEcho() (*echo.Echo, *MockPVZService, *PVZHandler) {
	e := echo.New()
	mockService := new(MockPVZService)
	s := &services.Services{}
	s.PvzService = mockService
	handler := NewPVZHandler(s)
	return e, mockService, handler
}

func TestPVZHandler_GetAll(t *testing.T) {
	e, mockService, handler := setupEcho()

	// Test case 1: Successful request
	t.Run("successful request", func(t *testing.T) {
		expectedPVZs := []models.FullPVZ{
			{
				ID:               "1",
				RegistrationDate: time.Now(),
				City:             "Moscow",
				Receptions:       make(map[string]models.FullReception),
			},
			{
				ID:               "2",
				RegistrationDate: time.Now(),
				City:             "Saint Petersburg",
				Receptions:       make(map[string]models.FullReception),
			},
		}
		mockService.On("GetAll", mock.Anything, "1", "10", "", "").
			Return(expectedPVZs, nil)

		req := httptest.NewRequest(http.MethodGet, "/pvz?page=1&limit=10", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.GetAll(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var response map[string]interface{}
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "1", response["page"])
		assert.Equal(t, "10", response["limit"])
	})

	// Test case 2: Error from service
	t.Run("service error", func(t *testing.T) {
		mockService.On("GetAll", mock.Anything, "", "", "", "").
			Return([]models.FullPVZ{}, assert.AnError)

		req := httptest.NewRequest(http.MethodGet, "/pvz", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.GetAll(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

func TestPVZHandler_Create(t *testing.T) {
	e, mockService, handler := setupEcho()

	// Test case 1: Successful creation
	t.Run("successful creation", func(t *testing.T) {
		pvz := models.PVZ{
			ID:               "1",
			RegistrationDate: time.Now(),
			City:             "Moscow",
		}
		mockService.On("CreatePVZ", mock.Anything, pvz).
			Return(pvz, nil)

		pvzJSON, _ := json.Marshal(pvz)
		req := httptest.NewRequest(http.MethodPost, "/pvz", strings.NewReader(string(pvzJSON)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.Create(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, rec.Code)
	})

	// Test case 2: Invalid request body
	t.Run("invalid request body", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/pvz", strings.NewReader("invalid json"))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.Create(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

func TestPVZHandler_GetByID(t *testing.T) {
	e, mockService, handler := setupEcho()

	// Test case 1: Successful retrieval
	t.Run("successful retrieval", func(t *testing.T) {
		expectedPVZ := models.PVZ{
			ID:               "1",
			RegistrationDate: time.Now(),
			City:             "Moscow",
		}
		mockService.On("GetPVZByID", mock.Anything, "1").
			Return(expectedPVZ, nil)

		req := httptest.NewRequest(http.MethodGet, "/pvz/1", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("1")

		err := handler.GetByID(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	// Test case 2: PVZ not found
	t.Run("pvz not found", func(t *testing.T) {
		mockService.On("GetPVZByID", mock.Anything, "2").
			Return(models.PVZ{}, assert.AnError)

		req := httptest.NewRequest(http.MethodGet, "/pvz/2", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("2")

		err := handler.GetByID(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})
}

func TestPVZHandler_DeleteLastProduct(t *testing.T) {
	e, mockService, handler := setupEcho()

	// Test case 1: Successful deletion
	t.Run("successful deletion", func(t *testing.T) {
		mockService.On("DeleteLastProduct", mock.Anything, "1").
			Return(nil)

		req := httptest.NewRequest(http.MethodPost, "/pvz/1/delete_last_product", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("1")

		err := handler.DeleteLastProduct(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	// Test case 2: PVZ not found
	t.Run("pvz not found", func(t *testing.T) {
		mockService.On("DeleteLastProduct", mock.Anything, "2").
			Return(assert.AnError)

		req := httptest.NewRequest(http.MethodPost, "/pvz/2/delete_last_product", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("2")

		err := handler.DeleteLastProduct(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})
}

func TestPVZHandler_CloseLastReception(t *testing.T) {
	e, mockService, handler := setupEcho()

	// Test case 1: Successful closure
	t.Run("successful closure", func(t *testing.T) {
		mockService.On("CloseLastReception", mock.Anything, "1").
			Return(nil)

		req := httptest.NewRequest(http.MethodPost, "/pvz/1/close_last_reception", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("1")

		err := handler.CloseLastReception(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	// Test case 2: PVZ not found
	t.Run("pvz not found", func(t *testing.T) {
		mockService.On("CloseLastReception", mock.Anything, "2").
			Return(assert.AnError)

		req := httptest.NewRequest(http.MethodPost, "/pvz/2/close_last_reception", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("2")

		err := handler.CloseLastReception(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})
}
