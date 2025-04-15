package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"pvz-service/config"
	"pvz-service/internal/handlers"
	"pvz-service/internal/models"
	"pvz-service/internal/repositories"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func setupTestPVZHandler(t *testing.T) (*handlers.PVZHandler, *echo.Echo) {
	cfg, err := config.NewConfig()
	assert.NoError(t, err)

	db, err := pgxpool.New(context.Background(), cfg.DATABASE_URL)
	assert.NoError(t, err)

	repo := repositories.NewPVZRepository(db)
	handler := handlers.NewPVZHandler(repo)

	e := echo.New()
	return handler, e
}

func TestPVZHandler_Create(t *testing.T) {
	handler, e := setupTestPVZHandler(t)
	validDate := time.Now().Format(time.RFC3339)
	validPayload := map[string]string{
		"city":             "Москва",
		"registrationDate": validDate,
	}
	validBody, _ := json.Marshal(validPayload)

	t.Run("successful creation", func(t *testing.T) {

		req := httptest.NewRequest(http.MethodPost, "/pvz", bytes.NewReader(validBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.Create(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, rec.Code)

		var resp models.PVZ
		err = json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, "москва", resp.City)
	})

	t.Run("unsupported city", func(t *testing.T) {
		payload := map[string]string{
			"city":             "Нижний Новгород",
			"registrationDate": validDate,
		}
		body, _ := json.Marshal(payload)

		req := httptest.NewRequest(http.MethodPost, "/pvz", bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.Create(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("invalid body", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/pvz", bytes.NewReader([]byte("invalid json")))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.Create(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}
