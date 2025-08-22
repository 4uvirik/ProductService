package handler

import (
	"github.com/4uvirik/ProductService/internal/entity"
	"github.com/4uvirik/ProductService/internal/usecase"
	"github.com/4uvirik/ProductService/pkg"
	"github.com/labstack/echo/v4"
	"log/slog"
	"net/http"
)

type ProductHandler struct {
	uc     *usecase.ProductUseCase
	logger *slog.Logger
}

func NewProductHandler(uc *usecase.ProductUseCase, logger *slog.Logger) *ProductHandler {
	return &ProductHandler{uc: uc, logger: logger}
}

func (h *ProductHandler) RegisterProductRoutes(g *echo.Group) {
	g.POST(entity.ProductURL, h.Create)
}

func (h *ProductHandler) Create(c echo.Context) error {
	var req entity.Product
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, pkg.ErrorResponse{Error: "invalid json"})
	}

	if err := h.uc.ProductCreate(c.Request().Context(), &req); err != nil {
		return c.JSON(pkg.MapErrorToStatus(err), pkg.ErrorResponse{Error: err.Error()})
	}
	return c.JSON(http.StatusCreated, req)
}
