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

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func setupTestReceptionHandler(t *testing.T) (*handlers.ReceptionHandler, *echo.Echo) {
	cfg, err := config.NewConfig()
	assert.NoError(t, err)

	db, err := pgxpool.New(context.Background(), cfg.DATABASE_URL)
	assert.NoError(t, err)

	repo := repositories.NewReceptionRepository(db)
	handler := handlers.NewReceptionHandler(repo)

	e := echo.New()
	return handler, e
}

func TestReceptionHandler_Create(t *testing.T) {
	handler, e := setupTestReceptionHandler(t)

	validPayload := map[string]string{
		"pvzId": "550e8400-e29b-41d4-a716-446655440000",
	}
	validBody, _ := json.Marshal(validPayload)

	t.Run("successful creation", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/receptions", bytes.NewReader(validBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.Create(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, rec.Code)

		var resp models.Reception
		err = json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, "in_progress", resp.Status)
	})

	t.Run("invalid body", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/receptions", bytes.NewReader([]byte("invalid json")))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.Create(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}
