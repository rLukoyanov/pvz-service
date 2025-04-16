package handlers

import (
	"net/http"
	"pvz-service/internal/pkg/jwt"
	"pvz-service/internal/services"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type DummyLoginHandler struct {
	services *services.Services
}

func NewDummyLoginHandler(services *services.Services) *DummyLoginHandler {
	return &DummyLoginHandler{
		services: services,
	}
}

// @Summary Получение тестового токена
// @Description Получение тестового токена для разработки
// @Tags auth
// @Accept json
// @Produce json
// @Param request body object true "User role data"
// @Param role body string true "User role" Enums(employee,moderator)
// @Success 200 {object} string "Token"
// @Failure 400 {object} map[string]string "Error message"
// @Router /dummyLogin [post]
func (h *DummyLoginHandler) DummyLogin(c echo.Context) error {
	type req struct {
		Role string `json:"role" enums:"employee,moderator"`
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

	token, err := jwt.GenerateToken(r.Role, h.services.Cfg)
	if err != nil {
		logrus.Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, "could not generate token")
	}

	return c.JSON(http.StatusOK, echo.Map{
		"token": token,
	})
}
