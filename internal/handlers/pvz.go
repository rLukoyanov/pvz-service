package handlers

import (
	"net/http"
	"pvz-service/internal/models"
	"pvz-service/internal/services"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type PVZHandler struct {
	services *services.Services
}

func NewPVZHandler(services *services.Services) *PVZHandler {
	return &PVZHandler{services: services}
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
		return echo.NewHTTPError(http.StatusBadRequest, echo.Map{"message": "invalid body"})
	}

	created, err := h.services.PvzService.CreatePVZ(c.Request().Context(), pvz)
	if err != nil {
		logrus.Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, echo.Map{"message": "could not create PVZ"})
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
	pvz, err := h.services.PvzService.GetPVZByID(c.Request().Context(), id)
	if err != nil {
		logrus.Error(err)
		return echo.NewHTTPError(http.StatusNotFound, echo.Map{"message": "PVZ not found"})
	}

	return c.JSON(http.StatusOK, pvz)
}

// @Summary Закрытие последней открытой приемки товаров в рамках ПВЗ
// @Description Закрытие последней открытой приемки товаров в рамках ПВЗ
// @Tags pvz
// @Security bearerAuth
// @Accept json
// @Produce json
// @Param id path string true "PVZ ID"
// @Success 200 {object} models.PVZ
// @Failure 404 {object} map[string]string
// @Router /pvz/{id}/delete_last_product [post]
func (h *PVZHandler) DeleteLastProduct(c echo.Context) error {
	id := c.Param("id")
	logrus.Info(id)
	err := h.services.PvzService.DeleteLastProduct(c.Request().Context(), id)
	if err != nil {
		logrus.Error(err)
		return echo.NewHTTPError(http.StatusNotFound, echo.Map{"message": "PVZ not found"})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "Товар удален"})
}

func (h *PVZHandler) CloseLastReception(c echo.Context) error {
	id := c.Param("id")
	logrus.Info(id)
	err := h.services.PvzService.CloseLastReception(c.Request().Context(), id)
	if err != nil {
		logrus.Error(err)
		return echo.NewHTTPError(http.StatusNotFound, echo.Map{"message": "PVZ not found"})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "Приемка закрыта"})
}
