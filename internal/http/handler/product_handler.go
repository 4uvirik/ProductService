package handler

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/4uvirik/ProductService/internal/entity"
	"github.com/4uvirik/ProductService/internal/http/httperr"
	"github.com/4uvirik/ProductService/internal/http/response"
	"github.com/4uvirik/ProductService/internal/usecase"
	"github.com/labstack/echo/v4"
)

type ProductHandler struct {
	uc     *usecase.ProductUseCase
	logger *slog.Logger
}

func NewProductHandler(uc *usecase.ProductUseCase, logger *slog.Logger) *ProductHandler {
	return &ProductHandler{uc: uc, logger: logger}
}

// RegisterProductRoutes - регистрирует маршруты для работы с продуктами.
func (h *ProductHandler) RegisterProductRoutes(g *echo.Group) {
	g.POST("/product", h.ProductCreate)
	g.GET("/product", h.ProductGetAll)
	g.GET("/product/:id", h.ProductGetByID)
	g.PUT("/product/:id", h.ProductUpdate)
	g.DELETE("/product/:id", h.ProductDelete)
}

// ProductCreate - обрабатывает POST запрос на создание нового продукта. Передает данные в usecase.
func (h *ProductHandler) ProductCreate(ctx echo.Context) error {
	var dto entity.ProductCreate
	if err := ctx.Bind(&dto); err != nil {
		h.logger.Error("bind error", slog.Any("err", err))

		return ctx.JSON(httperr.MapErrorToStatus(err), httperr.ErrorResponse{Error: "invalid json"})
	}

	product, err := h.uc.ProductCreate(ctx.Request().Context(), dto)
	if err != nil {
		return ctx.JSON(httperr.MapErrorToStatus(err), httperr.ErrorResponse{Error: err.Error()})
	}

	return ctx.JSON(http.StatusCreated, product)
}

// ProductGetAll - обрабатывает GET запрос для получения всех продуктов.
func (h *ProductHandler) ProductGetAll(ctx echo.Context) error {
	allProducts, err := h.uc.ProductGetAll(ctx.Request().Context())
	if err != nil {
		return ctx.JSON(httperr.MapErrorToStatus(err), httperr.ErrorResponse{Error: err.Error()})
	}

	return ctx.JSON(http.StatusOK, allProducts)
}

// ProductGetByID - обрабатывает GET запрос на получение продукта по id. Запрашивает продукт из usecase.
func (h *ProductHandler) ProductGetByID(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(httperr.MapErrorToStatus(err), httperr.ErrorResponse{Error: "invalid id"})
	}

	product, err := h.uc.ProductGetByID(ctx.Request().Context(), id)
	if err != nil {
		return ctx.JSON(httperr.MapErrorToStatus(err), httperr.ErrorResponse{Error: err.Error()})
	}

	return ctx.JSON(http.StatusOK, product)
}

// ProductUpdate - обрабатывает PUT запрос на обновление продукта. Декодирует тело запроса и обновляет продукт.
func (h *ProductHandler) ProductUpdate(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(httperr.MapErrorToStatus(err), httperr.ErrorResponse{Error: "invalid id"})
	}

	var dto entity.ProductUpdate
	if err := ctx.Bind(&dto); err != nil {
		return ctx.JSON(httperr.MapErrorToStatus(err), httperr.ErrorResponse{Error: "invalid request body"})
	}

	dto.ID = id

	if err := h.uc.ProductUpdate(ctx.Request().Context(), &dto); err != nil {
		return ctx.JSON(httperr.MapErrorToStatus(err), httperr.ErrorResponse{Error: err.Error()})
	}

	return ctx.JSON(http.StatusOK, response.Response{Result: "updated"})
}

// ProductDelete - обрабатывает DELETE запрос на удаление продукта по id.
func (h *ProductHandler) ProductDelete(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(httperr.MapErrorToStatus(err), httperr.ErrorResponse{Error: "invalid id"})
	}

	if err := h.uc.ProductDelete(ctx.Request().Context(), id); err != nil {
		return ctx.JSON(httperr.MapErrorToStatus(err), httperr.ErrorResponse{Error: err.Error()})
	}

	return ctx.JSON(http.StatusOK, response.Response{Result: "deleted"})
}
