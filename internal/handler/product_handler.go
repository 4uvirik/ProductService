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

// RegisterProductRoutes - регистрирует маршруты для работы с продуктами.
// Настраивает пути и методы:
// POST   /products         - создание продукта
// GET    /products         - получение списка всех продуктов
// GET    /products/:id     - получение продукта по id
// PUT    /products/:id     - обновление данных продукта по id
// DELETE /products/:id     - удаление продукта по id
func (h *ProductHandler) RegisterProductRoutes(g *echo.Group) {
	g.POST(entity.ProductURL, h.ProductCreate)
	g.GET(entity.ProductURL, h.ProductGetAll)
	g.GET(entity.ProductURL+entity.IDParam, h.ProductGetByID)
	g.PUT(entity.ProductURL+entity.IDParam, h.ProductUpdate)
	g.DELETE(entity.ProductURL+entity.IDParam, h.ProductDelete)
}

// ProductCreate - обрабатывает POST запрос на создание нового продукта. Передает данные в usecase
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

// ProductGetAll - обрабатывает GET запрос для получения всех продуктов
func (h *ProductHandler) ProductGetAll(c echo.Context) error {
	allProducts, err := h.uc.ProductGetAll(c.Request().Context())
	if err != nil {
		return c.JSON(pkg.MapErrorToStatus(err), pkg.ErrorResponse{Error: err.Error()})
	}
	return c.JSON(http.StatusOK, allProducts)
}

// ProductGetByID - обрабатывает GET запрос на получение продукта по id. Запрашивает продукт из usecase
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

// ProductUpdate - обрабатывает PUT запрос на обновление продукта. Декодирует тело запроса и обновляет продукт
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

// ProductDelete - обрабатывает DELETE запрос на удаление продукта по id
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
