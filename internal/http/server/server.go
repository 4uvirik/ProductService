package server

import (
	"github.com/4uvirik/ProductService/config"
	"github.com/4uvirik/ProductService/internal/http/handler"
	"github.com/4uvirik/ProductService/internal/usecase"
	"github.com/labstack/echo/v4"
	"log/slog"
)

type Server struct {
	cfg        config.AppConfig
	logger     *slog.Logger
	UCCategory *usecase.CategoryUseCase
	UCProduct  *usecase.ProductUseCase
}

func Run(cfg *config.Config, UCCategory *usecase.CategoryUseCase, UCProduct *usecase.ProductUseCase, logger *slog.Logger) {
	e := echo.New()

	handlerCategory := handler.NewCategoryHandler(UCCategory, logger)
	handlerProduct := handler.NewProductHandler(UCProduct, logger)

	api := e.Group("")
	handlerCategory.RegisterCategoryRoutes(api)
	handlerProduct.RegisterProductRoutes(api)

	address := cfg.App.Host + ":" + cfg.App.Port
	logger.Info("server started", slog.String("address", address))
	e.Logger.Fatal(e.Start(address))
}
