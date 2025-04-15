package handlers

import (
	"net/http"
	"pvz-service/internal/models"
	"pvz-service/internal/repositories"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type PVZHandler struct {
	repo *repositories.PVZRepository
}

func NewPVZHandler(repo *repositories.PVZRepository) *PVZHandler {
	return &PVZHandler{repo: repo}
}

func (h *PVZHandler) Create(c echo.Context) error {
	var req models.PVZ
	if err := c.Bind(&req); err != nil {
		logrus.Error(err)
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "invalid body"})
	}
	req.City = strings.ToLower(req.City)

	allowedCities := map[string]struct{}{
		"москва":          {},
		"санкт-Петербург": {},
		"казань":          {},
	}

	if _, ok := allowedCities[req.City]; !ok {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "unsupported city"})
	}

	point := models.PVZ{
		RegistrationDate: req.RegistrationDate,
		City:             req.City,
	}

	created, err := h.repo.CreatePVZ(c.Request().Context(), point)
	if err != nil {
		logrus.Error(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "failed to create pvz"})
	}

	return c.JSON(http.StatusCreated, created)
}
