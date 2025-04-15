package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"pvz-service/config"
	"pvz-service/internal/handlers"

	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestDummyLogin(t *testing.T) {
	e := echo.New()

	cfg, err := config.NewConfig()
	if err != nil {
		panic(err)
	}

	dlHandler := handlers.NewDummyLoginHandler(cfg)
	e.POST("/dummyLogin", dlHandler.DummyLogin)

	t.Run("valid role - employee", func(t *testing.T) {
		body := `{"role":"client"}`
		req := httptest.NewRequest(http.MethodPost, "/dummyLogin", bytes.NewBufferString(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)

		var response map[string]string
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotEmpty(t, response["token"])
	})

	t.Run("invalid role", func(t *testing.T) {
		body := `{"role":"admin"}`
		req := httptest.NewRequest(http.MethodPost, "/dummyLogin", bytes.NewBufferString(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("empty body", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/dummyLogin", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}
