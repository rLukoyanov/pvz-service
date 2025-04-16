package handlers

import (
	"net/http"
	"pvz-service/internal/models"
	"pvz-service/internal/repositories"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type ReceptionHandler struct {
	Repo *repositories.ReceptionRepository
}

func NewReceptionHandler(repo *repositories.ReceptionRepository) *ReceptionHandler {
	return &ReceptionHandler{Repo: repo}
}

// @Summary Создание новой приемки товаров
// @Description Создание новой приемки товаров (только для сотрудников ПВЗ)
// @Tags Reception
// @Security bearerAuth
// @Accept json
// @Produce json
// @Param request body models.Reception true "Reception data"
// @Success 201 {object} models.Reception
// @Failure 400 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Router /receptions [post]
func (h *ReceptionHandler) Create(c echo.Context) error {
	var req struct {
		PvzId string `json:"pvzId"`
	}
	if err := c.Bind(&req); err != nil {
		logrus.Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, map[string]string{"message": "invalid body"})
	}

	active, err := h.Repo.GetActiveReceptionByPVZID(c.Request().Context(), req.PvzId)
	if err != nil {
		logrus.Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, map[string]string{"message": "could not create Reception"})
	}

	if active != nil {
		return echo.NewHTTPError(http.StatusBadRequest, map[string]string{"message": "there is an active Reception for this PVZ"})
	}

	Reception := models.Reception{
		PvzId:  req.PvzId,
		Status: "in_progress",
	}

	err = h.Repo.CreateReception(c.Request().Context(), Reception)
	if err != nil {
		logrus.Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, map[string]string{"message": "could not create Reception"})
	}

	return c.JSON(http.StatusCreated, map[string]string{"message": "created"})
}
