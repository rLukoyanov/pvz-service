package handlers

import (
	"net/http"
	"pvz-service/config"
	"pvz-service/internal/pkg/jwt"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type DummyLoginHandler struct {
	cfg *config.Config
}

func NewDummyLoginHandler(cfg *config.Config) *DummyLoginHandler {
	return &DummyLoginHandler{
		cfg: cfg,
	}
}

func (h *DummyLoginHandler) DummyLogin(c echo.Context) error {
	type req struct {
		Role string `json:"role"`
	}

	var r req
	if err := c.Bind(&r); err != nil {
		logrus.Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, "invalid body")
	}

	if r.Role != "client" && r.Role != "moderator" {
		logrus.Error("invalid role")
		return echo.NewHTTPError(http.StatusBadRequest, "invalid role")
	}

	token, err := jwt.GenerateToken(r.Role, h.cfg)
	if err != nil {
		logrus.Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, "could not generate token")
	}

	return c.JSON(http.StatusOK, echo.Map{
		"token": token,
	})
}
