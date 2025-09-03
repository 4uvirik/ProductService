package main

import (
	"context"
	"log"
	"log/slog"
	"os"

	"github.com/4uvirik/ProductService/config"
	"github.com/4uvirik/ProductService/internal/http/server"
	"github.com/4uvirik/ProductService/internal/repository"
	"github.com/4uvirik/ProductService/internal/usecase"
	"github.com/4uvirik/ProductService/pkg/logger"
	"github.com/joho/godotenv"
)

func main() {
	loadEnv()

	cfg := initConfig()

	logger := initLogger(cfg)

	postgres := initPostgres(cfg, logger)

	defer postgres.Pool.Close()

	categoryUC, productUC := initUseCase(postgres, logger)

	server.Run(cfg, categoryUC, productUC, logger)
}

// loadEnv - загрузка переменных окружения.
func loadEnv() {
	if err := godotenv.Load(".env"); err != nil {
		log.Println(".env file not found")
	}
}

// initConfig - загрузка конфига и проверка на валидацию.
func initConfig() *config.Config {
	yamlPath := os.Getenv("YAML_PATH")

	cfg, err := config.LoadConfig(yamlPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	if err := cfg.Validate(); err != nil {
		log.Fatalf("config validate failed: %v", err)
	}

	return cfg
}

// initLogger - загрузка и настройка логера.
func initLogger(cfg *config.Config) *slog.Logger {
	logger, err := logger.ConfigLogger(cfg.Logger.Level)
	if err != nil {
		log.Fatalf("error creating logger: %v", err)
	}

	logger.Debug("configuration reading success", slog.Any("cfg", cfg))

	return logger
}

// initPostgres - подключение к базе данных по DSN из конфига.
func initPostgres(cfg *config.Config, logger *slog.Logger) *repository.Postgres {
	postgres, err := repository.NewPostgres(context.Background(), cfg.DSN(), logger)
	if err != nil {
		log.Fatalf("error creating new repository: %v", err)
	}

	return postgres
}

// initUseCase - загрузка репозитория и бизнес логики (usecase).
func initUseCase(postgres *repository.Postgres, logger *slog.Logger) (*usecase.CategoryUseCase, *usecase.ProductUseCase) {
	repoCategory := repository.NewCategoryRepository(postgres.Pool, logger)
	repoProduct := repository.NewProductRepository(postgres.Pool, logger)

	categoryUC := usecase.NewCategoryUseCase(repoCategory, logger)
	productUC := usecase.NewProductUseCase(repoProduct, logger)

	return categoryUC, productUC
}
