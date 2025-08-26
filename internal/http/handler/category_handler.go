package handler

import (
	"encoding/json"
	"github.com/4uvirik/ProductService/internal/entity"
	"github.com/4uvirik/ProductService/internal/http/httperr"
	"github.com/4uvirik/ProductService/internal/http/response"
	"github.com/4uvirik/ProductService/internal/usecase"
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

// RegisterCategoryRoutes - регистрирует маршруты для работы с категориями.
// Настраивает пути и методы:
// POST   /category         - создание продукта
// PUT    /category/:id     - обновление данных продукта по id
// DELETE /category/:id     - удаление продукта по id
func (h *CategoryHandler) RegisterCategoryRoutes(g *echo.Group) {
	g.POST(entity.CategoryURL, h.CategoryCreate)
	g.PUT(entity.CategoryURL+entity.IDParam, h.CategoryUpdate)
	g.DELETE(entity.CategoryURL+entity.IDParam, h.CategoryDelete)
}

// CategoryCreate - обрабатывает POST запрос на создание новой категории продуктов. Передает данные в usecase
func (h *CategoryHandler) CategoryCreate(c echo.Context) error {
	var req entity.Category
	if err := c.Bind(&req); err != nil {
		return c.JSON(httperr.MapErrorToStatus(err), httperr.ErrorResponse{Error: "invalid json"})
	}

	if err := h.uc.CategoryCreate(c.Request().Context(), &req); err != nil {
		return c.JSON(httperr.MapErrorToStatus(err), httperr.ErrorResponse{Error: err.Error()})
	}
	return c.JSON(http.StatusCreated, req)
}

// CategoryUpdate - обрабатывает PUT запрос на обновление названия категории продуктов. Декодирует тело запроса и обновляет продукт
func (h *CategoryHandler) CategoryUpdate(c echo.Context) error {
	var upd entity.Category
	if err := json.NewDecoder(c.Request().Body).Decode(&upd); err != nil {
		return c.JSON(httperr.MapErrorToStatus(err), httperr.ErrorResponse{Error: "invalid request body"})
	}

	if err := h.uc.CategoryUpdate(c.Request().Context(), &upd); err != nil {
		return c.JSON(httperr.MapErrorToStatus(err), httperr.ErrorResponse{Error: err.Error()})
	}
	return c.JSON(http.StatusOK, response.Response{Result: "updated"})
}

// CategoryDelete - обрабатывает DELETE запрос на удаление категории продуктов по id
func (h *CategoryHandler) CategoryDelete(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(httperr.MapErrorToStatus(err), httperr.ErrorResponse{Error: "invalid id"})
	}

	if err := h.uc.CategoryDelete(c.Request().Context(), id); err != nil {
		return c.JSON(httperr.MapErrorToStatus(err), httperr.ErrorResponse{Error: err.Error()})
	}
	return c.JSON(http.StatusOK, response.Response{Result: "deleted"})
}
