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

//nolint:cyclop
func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Println(".env file not found")
	}

	yamlPath := os.Getenv("YAML_PATH")

	cfg, err := config.LoadConfig(yamlPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	if err := cfg.Validate(); err != nil {
		log.Fatalf("config validate failed: %v", err)
	}

	logger, err := logger.ConfigLogger(cfg.Logger.Level)
	if err != nil {
		log.Fatalf("error creating logger: %v", err)
	}

	logger.Debug("configuration reading success", slog.Any("cfg", cfg))

	postgres, err := repository.NewPostgres(context.Background(), cfg.DSN(), logger)
	if err != nil {
		log.Fatalf("error creating new repository: %v", err)
	}
	defer postgres.Pool.Close()

	repoCategory := repository.NewCategoryRepository(postgres.Pool, logger)
	repoProduct := repository.NewProductRepository(postgres.Pool, logger)

	categoryUC := usecase.NewCategoryUseCase(repoCategory, logger)
	productUC := usecase.NewProductUseCase(repoProduct, logger)

	server.Run(cfg, categoryUC, productUC, logger)
}
