package handlers

import (
	"net/http"
	"pvz-service/internal/models"
	"pvz-service/internal/services"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type ItemHandler struct {
	services *services.Services
}

func NewProductHandler(services *services.Services) *ItemHandler {
	return &ItemHandler{services: services}
}

type req struct {
	Type  string `json:"type"`
	PvzId string `json:"PvzId"`
}

func (h *ItemHandler) AddProduct(c echo.Context) error {
	var req req
	if err := c.Bind(&req); err != nil {
		logrus.Error("Failed to bind request:", err)
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	if req.Type == "" || req.PvzId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "missing required fields")
	}

	product := models.Product{
		Type: strings.ToLower(req.Type),
	}

	if err := h.services.ProductService.AddProduct(c.Request().Context(), product, req.PvzId); err != nil {
		logrus.Error("Failed to add product:", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to add product")
	}

	return c.JSON(http.StatusCreated, echo.Map{"message": "created"})
}
