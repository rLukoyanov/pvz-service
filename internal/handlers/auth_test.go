package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"pvz-service/config"
	"pvz-service/internal/handlers"
	"pvz-service/internal/repositories"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

func setupTestAuthHandler(t *testing.T) (*handlers.AuthHandler, *echo.Echo) {
	cfg, err := config.NewConfig()
	assert.NoError(t, err)

	db, err := pgxpool.New(context.Background(), cfg.DATABASE_URL)
	assert.NoError(t, err)

	repo := repositories.NewUserRepository(db)
	handler := handlers.NewAuthHandler(repo, cfg)

	e := echo.New()
	return handler, e
}

func TestRegisterLogin(t *testing.T) {
	handler, e := setupTestAuthHandler(t)

	email := "testuser@example.com"
	password := "securePassword"
	role := "client"
	t.Run("register user", func(t *testing.T) {
		reqBody := map[string]string{
			"email":    email,
			"password": password,
			"role":     role,
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.Register(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.Contains(t, rec.Body.String(), email)
	})

	t.Run("login with correct credentials", func(t *testing.T) {
		reqBody := map[string]string{
			"email":    email,
			"password": password,
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.Login(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "token")
	})

	t.Run("login with wrong password", func(t *testing.T) {
		reqBody := map[string]string{
			"email":    email,
			"password": "wrongPassword",
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.Login(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})
}
