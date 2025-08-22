package handler

import (
	"encoding/json"
	"github.com/4uvirik/ProductService/internal/entity"
	"github.com/4uvirik/ProductService/internal/usecase"
	"github.com/4uvirik/ProductService/pkg"
	"github.com/labstack/echo/v4"
	"log/slog"
	"net/http"
	"strconv"
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
	g.GET(entity.ProductURL, h.ProductGetAll)
	g.GET(entity.ProductURL+entity.IDParam, h.ProductGetByID)
	g.PUT(entity.ProductURL+entity.IDParam, h.ProductUpdate)
	g.DELETE(entity.ProductURL+entity.IDParam, h.ProductDelete)
}

func (h *ProductHandler) ProductCreate(c echo.Context) error {
	var req entity.Product
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, pkg.ErrorResponse{Error: "invalid json"})
	}

	if err := h.uc.ProductCreate(c.Request().Context(), &req); err != nil {
		return c.JSON(pkg.MapErrorToStatus(err), pkg.ErrorResponse{Error: err.Error()})
	}
	return c.JSON(http.StatusCreated, req)
}

func (h *ProductHandler) ProductGetAll(c echo.Context) error {
	allProducts, err := h.uc.ProductGetAll(c.Request().Context())
	if err != nil {
		return c.JSON(pkg.MapErrorToStatus(err), pkg.ErrorResponse{Error: err.Error()})
	}
	return c.JSON(http.StatusOK, allProducts)
}

func (h *ProductHandler) ProductGetByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, pkg.ErrorResponse{Error: "invalid id"})
	}

	product, err := h.uc.ProductGetByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(pkg.MapErrorToStatus(err), pkg.ErrorResponse{Error: err.Error()})
	}
	return c.JSON(http.StatusOK, product)
}

func (h *ProductHandler) ProductUpdate(c echo.Context) error {
	var upd entity.ProductUpdate
	if err := json.NewDecoder(c.Request().Body).Decode(&upd); err != nil {
		return c.JSON(http.StatusBadRequest, pkg.ErrorResponse{Error: "invalid request body"})
	}

	if err := h.uc.ProductUpdate(c.Request().Context(), &upd); err != nil {
		return c.JSON(http.StatusBadRequest, pkg.ErrorResponse{Error: err.Error()})
	}
	return c.JSON(http.StatusOK, pkg.Response{Result: "updated"})
}

func (h *ProductHandler) ProductDelete(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, pkg.ErrorResponse{Error: "invalid id"})
	}

	if err := h.uc.ProductDelete(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusBadRequest, pkg.ErrorResponse{Error: err.Error()})
	}
	return c.JSON(http.StatusOK, pkg.Response{Result: "deleted"})
}
