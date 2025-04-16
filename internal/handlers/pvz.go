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
	Repo *repositories.PVZRepository
}

func NewPVZHandler(repo *repositories.PVZRepository) *PVZHandler {
	return &PVZHandler{Repo: repo}
}

// @Summary Создание ПВЗ
// @Description Создание нового пункта выдачи заказов (только для модераторов)
// @Tags pvz
// @Security bearerAuth
// @Accept json
// @Produce json
// @Param request body models.PVZ true "PVZ data"
// @Success 201 {object} models.PVZ
// @Failure 400 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Router /pvz [post]
func (h *PVZHandler) Create(c echo.Context) error {
	var pvz models.PVZ
	if err := c.Bind(&pvz); err != nil {
		logrus.Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, map[string]string{"message": "invalid body"})
	}

	pvz.City = strings.ToLower(pvz.City)

	allowedCities := map[string]bool{
		"москва":          true,
		"санкт-Петербург": true,
		"казань":          true,
	}

	if _, ok := allowedCities[pvz.City]; !ok {
		return echo.NewHTTPError(http.StatusBadRequest, map[string]string{"message": "city not allowed"})
	}

	created, err := h.Repo.CreatePVZ(c.Request().Context(), pvz)
	if err != nil {
		logrus.Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, map[string]string{"message": "could not create PVZ"})
	}

	return c.JSON(http.StatusCreated, created)
}

// @Summary Получение ПВЗ по ID
// @Description Получение информации о пункте выдачи заказов по его ID
// @Tags pvz
// @Security bearerAuth
// @Accept json
// @Produce json
// @Param id path string true "PVZ ID"
// @Success 200 {object} models.PVZ
// @Failure 404 {object} map[string]string
// @Router /pvz/{id} [get]
func (h *PVZHandler) GetByID(c echo.Context) error {
	id := c.Param("id")
	pvz, err := h.Repo.GetPVZByID(c.Request().Context(), id)
	if err != nil {
		logrus.Error(err)
		return echo.NewHTTPError(http.StatusNotFound, map[string]string{"message": "PVZ not found"})
	}

	return c.JSON(http.StatusOK, pvz)
}
