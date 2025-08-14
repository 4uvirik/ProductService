package main

import (
	"github.com/4uvirik/ProductService/config"
	"github.com/4uvirik/ProductService/internal/logger"
	"log"
	"log/slog"
	"os"
)

func main() {

	yamlPath := os.Getenv("YAML_PATH")
	if yamlPath == "" {
		yamlPath = "config/config.yaml"
	}

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

}
