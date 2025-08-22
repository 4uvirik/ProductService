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

type CategoryHandler struct {
	uc     *usecase.CategoryUseCase
	logger *slog.Logger
}

func NewCategoryHandler(uc *usecase.ProductUseCase, logger *slog.Logger) *ProductHandler {
	return &ProductHandler{uc: uc, logger: logger}
}

func (h *CategoryHandler) RegisterCategoryRoutes(g *echo.Group) {
	g.POST(entity.ProductURL, h.CategoryCreate)
	g.PUT(entity.ProductURL+entity.IDParam, h.CategoryUpdate)
	g.DELETE(entity.ProductURL+entity.IDParam, h.CategoryDelete)
}

func (h *CategoryHandler) CategoryCreate(c echo.Context) error {
	var req entity.Category
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, pkg.ErrorResponse{Error: "invalid json"})
	}

	if err := h.uc.CategoryCreate(c.Request().Context(), &req); err != nil {
		return c.JSON(pkg.MapErrorToStatus(err), pkg.ErrorResponse{Error: err.Error()})
	}
	return c.JSON(http.StatusCreated, req)
}

func (h *CategoryHandler) CategoryUpdate(c echo.Context) error {
	var upd entity.Category
	if err := json.NewDecoder(c.Request().Body).Decode(&upd); err != nil {
		return c.JSON(http.StatusBadRequest, pkg.ErrorResponse{Error: "invalid request body"})
	}

	if err := h.uc.CategoryUpdate(c.Request().Context(), &upd); err != nil {
		return c.JSON(http.StatusBadRequest, pkg.ErrorResponse{Error: err.Error()})
	}
	return c.JSON(http.StatusOK, pkg.Response{Result: "updated"})
}

func (h *CategoryHandler) CategoryDelete(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, pkg.ErrorResponse{Error: "invalid id"})
	}

	if err := h.uc.CategoryDelete(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusBadRequest, pkg.ErrorResponse{Error: err.Error()})
	}
	return c.JSON(http.StatusOK, pkg.Response{Result: "deleted"})
}
