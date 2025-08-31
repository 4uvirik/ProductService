package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/4uvirik/ProductService/internal/entity"
	"github.com/4uvirik/ProductService/internal/http/httperr"
	"github.com/4uvirik/ProductService/internal/http/response"
	"github.com/4uvirik/ProductService/internal/usecase"
	"github.com/labstack/echo/v4"
)

type CategoryHandler struct {
	uc     *usecase.CategoryUseCase
	logger *slog.Logger
}

func NewCategoryHandler(uc *usecase.CategoryUseCase, logger *slog.Logger) *CategoryHandler {
	return &CategoryHandler{uc: uc, logger: logger}
}

// RegisterCategoryRoutes - регистрирует маршруты для работы с категориями.
func (h *CategoryHandler) RegisterCategoryRoutes(g *echo.Group) {
	g.POST(entity.CategoryURL, h.CategoryCreate)
	g.PUT(entity.CategoryURL+entity.IDParam, h.CategoryUpdate)
	g.DELETE(entity.CategoryURL+entity.IDParam, h.CategoryDelete)
}

// CategoryCreate - обрабатывает POST запрос на создание новой категории продуктов. Передает данные в usecase.
func (h *CategoryHandler) CategoryCreate(ctx echo.Context) error {
	var req entity.Category
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(httperr.MapErrorToStatus(err), httperr.ErrorResponse{Error: "invalid json"})
	}

	if err := h.uc.CategoryCreate(ctx.Request().Context(), &req); err != nil {
		return ctx.JSON(httperr.MapErrorToStatus(err), httperr.ErrorResponse{Error: err.Error()})
	}

	return ctx.JSON(http.StatusCreated, req)
}

// CategoryUpdate - обрабатывает PUT запрос на обновление названия категории продуктов. Декодирует тело запроса и обновляет продукт.
func (h *CategoryHandler) CategoryUpdate(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(httperr.MapErrorToStatus(err), httperr.ErrorResponse{Error: "invalid id"})
	}

	var upd entity.Category
	if err := json.NewDecoder(ctx.Request().Body).Decode(&upd); err != nil {
		return ctx.JSON(httperr.MapErrorToStatus(err), httperr.ErrorResponse{Error: "invalid request body"})
	}

	upd.ID = id

	if err := h.uc.CategoryUpdate(ctx.Request().Context(), &upd); err != nil {
		return ctx.JSON(httperr.MapErrorToStatus(err), httperr.ErrorResponse{Error: err.Error()})
	}

	return ctx.JSON(http.StatusOK, response.Response{Result: "updated"})
}

// CategoryDelete - обрабатывает DELETE запрос на удаление категории продуктов по id.
func (h *CategoryHandler) CategoryDelete(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(httperr.MapErrorToStatus(err), httperr.ErrorResponse{Error: "invalid id"})
	}

	if err := h.uc.CategoryDelete(ctx.Request().Context(), id); err != nil {
		return ctx.JSON(httperr.MapErrorToStatus(err), httperr.ErrorResponse{Error: err.Error()})
	}

	return ctx.JSON(http.StatusOK, response.Response{Result: "deleted"})
}
