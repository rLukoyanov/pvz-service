package handlers

import (
	"net/http"
	"pvz-service/internal/models"
	"pvz-service/internal/repositories"

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
		return echo.NewHTTPError(http.StatusBadRequest, "invalid body")
	}

	allowedCities := map[string]bool{
		"москва":          true,
		"санкт-петербург": true,
		"казань":          true,
	}

	if !allowedCities[pvz.City] {
		return echo.NewHTTPError(http.StatusBadRequest, "city not allowed")
	}

	created, err := h.Repo.CreatePVZ(c.Request().Context(), pvz)
	if err != nil {
		logrus.Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, "could not create PVZ")
	}

	return c.JSON(http.StatusCreated, created)
}

// @Summary Получение ПВЗ по ID
// @Description Получение информации о пункте выдачи заказов по его ID
// @Tags pvz
// @Security bearerAuth
// @Accept json
// @Produce json
// @Param id path int true "PVZ ID"
// @Success 200 {object} models.PVZ
// @Failure 404 {object} map[string]string
// @Router /pvz/{id} [get]
func (h *PVZHandler) GetByID(c echo.Context) error {
	id := c.Param("id")
	pvz, err := h.Repo.GetPVZByID(c.Request().Context(), id)
	if err != nil {
		logrus.Error(err)
		return echo.NewHTTPError(http.StatusNotFound, "PVZ not found")
	}

	return c.JSON(http.StatusOK, pvz)
}

// @Summary Обновление ПВЗ
// @Description Обновление информации о пункте выдачи заказов (только для модераторов)
// @Tags pvz
// @Security bearerAuth
// @Accept json
// @Produce json
// @Param id path int true "PVZ ID"
// @Param request body models.PVZ true "PVZ data"
// @Success 200 {object} models.PVZ
// @Failure 400 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /pvz/{id} [put]
func (h *PVZHandler) Update(c echo.Context) error {
	id := c.Param("id")
	var pvz models.PVZ
	if err := c.Bind(&pvz); err != nil {
		logrus.Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, "invalid body")
	}

	allowedCities := map[string]bool{
		"москва":          true,
		"санкт-петербург": true,
		"казань":          true,
	}

	if !allowedCities[pvz.City] {
		return echo.NewHTTPError(http.StatusBadRequest, "city not allowed")
	}

	updated, err := h.Repo.UpdatePVZ(c.Request().Context(), id, pvz)
	if err != nil {
		logrus.Error(err)
		return echo.NewHTTPError(http.StatusNotFound, "PVZ not found")
	}

	return c.JSON(http.StatusOK, updated)
}

// @Summary Удаление ПВЗ
// @Description Удаление пункта выдачи заказов (только для модераторов)
// @Tags pvz
// @Security bearerAuth
// @Accept json
// @Produce json
// @Param id path int true "PVZ ID"
// @Success 204 "No Content"
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /pvz/{id} [delete]
func (h *PVZHandler) Delete(c echo.Context) error {
	id := c.Param("id")
	if err := h.Repo.DeletePVZ(c.Request().Context(), id); err != nil {
		logrus.Error(err)
		return echo.NewHTTPError(http.StatusNotFound, "PVZ not found")
	}

	return c.NoContent(http.StatusNoContent)
}
