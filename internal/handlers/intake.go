package handlers

import (
	"net/http"
	"pvz-service/internal/models"
	"pvz-service/internal/repositories"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type IntakeHandler struct {
	Repo *repositories.IntakeRepository
}

func NewIntakeHandler(repo *repositories.IntakeRepository) *IntakeHandler {
	return &IntakeHandler{Repo: repo}
}

// @Summary Создание новой приемки товаров
// @Description Создание новой приемки товаров (только для сотрудников ПВЗ)
// @Tags intake
// @Security bearerAuth
// @Accept json
// @Produce json
// @Param request body object true "Intake data"
// @Param pvzId body string true "PVZ ID" Format(uuid)
// @Success 201 {object} models.Intake
// @Failure 400 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Router /receptions [post]
func (h *IntakeHandler) Create(c echo.Context) error {
	var req struct {
		PvzId string `json:"pvzId"`
	}
	if err := c.Bind(&req); err != nil {
		logrus.Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, map[string]string{"message": "invalid body"})
	}

	active, err := h.Repo.GetActiveIntakeByPVZID(c.Request().Context(), req.PvzId)
	if err != nil {
		logrus.Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, map[string]string{"message": "could not create intake"})
	}

	if active != nil {
		return echo.NewHTTPError(http.StatusBadRequest, map[string]string{"message": "there is an active intake for this PVZ"})
	}

	intake := models.Intake{
		PvzId:  req.PvzId,
		Status: "in_progress",
	}

	err = h.Repo.CreateIntake(c.Request().Context(), intake)
	if err != nil {
		logrus.Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, map[string]string{"message": "could not create intake"})
	}

	return c.JSON(http.StatusCreated, map[string]string{"message": "created"})
}
