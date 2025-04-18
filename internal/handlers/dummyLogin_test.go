package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"pvz-service/config"
	"pvz-service/internal/services"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func setupDummyLoginEcho() (*echo.Echo, *DummyLoginHandler) {
	e := echo.New()
	s := &services.Services{
		Cfg: &config.Config{
			SECRET: "test_secret",
		},
	}
	handler := NewDummyLoginHandler(s)
	return e, handler
}

func TestDummyLoginHandler_DummyLogin(t *testing.T) {
	e, handler := setupDummyLoginEcho()

	// Test case 1: Successful token generation for client
	t.Run("successful token generation for client", func(t *testing.T) {
		reqBody := map[string]string{
			"role": "client",
		}
		reqJSON, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/dummyLogin", strings.NewReader(string(reqJSON)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.DummyLogin(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var response map[string]string
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotEmpty(t, response["token"])
	})

	// Test case 2: Successful token generation for moderator
	t.Run("successful token generation for moderator", func(t *testing.T) {
		reqBody := map[string]string{
			"role": "moderator",
		}
		reqJSON, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/dummyLogin", strings.NewReader(string(reqJSON)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.DummyLogin(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var response map[string]string
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotEmpty(t, response["token"])
	})

	// Test case 3: Invalid request body
	t.Run("invalid request body", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/dummyLogin", strings.NewReader("invalid json"))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.DummyLogin(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	// Test case 4: Invalid role
	t.Run("invalid role", func(t *testing.T) {
		reqBody := map[string]string{
			"role": "invalid_role",
		}
		reqJSON, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/dummyLogin", strings.NewReader(string(reqJSON)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.DummyLogin(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}
