package server

import (
	"log/slog"

	"github.com/4uvirik/ProductService/config"
	"github.com/4uvirik/ProductService/internal/http/handler"
	"github.com/4uvirik/ProductService/internal/usecase"
	"github.com/labstack/echo/v4"
)

func Run(cfg *config.Config, categoryUC *usecase.CategoryUseCase, productUC *usecase.ProductUseCase, logger *slog.Logger) *echo.Echo {
	e := echo.New()

	handlerCategory := handler.NewCategoryHandler(categoryUC, logger)
	handlerProduct := handler.NewProductHandler(productUC, logger)

	api := e.Group("")
	handlerCategory.RegisterCategoryRoutes(api)
	handlerProduct.RegisterProductRoutes(api)

	address := cfg.App.Host + ":" + cfg.App.Port
	logger.Info("server started", slog.String("address", address))

	return e
}
