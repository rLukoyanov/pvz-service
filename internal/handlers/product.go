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
	req := req{}
	if err := c.Bind(&req); err != nil || req.Type == "" || req.PvzId == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"message": "invalid request"})
	}

	product := models.Product{}
	product.Type = strings.ToLower(req.Type)

	if err := h.services.ProductService.AddProduct(c.Request().Context(), product, req.PvzId); err != nil {
		logrus.Error(err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "failed to add item"})
	}

	return c.JSON(http.StatusCreated, echo.Map{"message": "created"})
}
