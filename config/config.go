package config

import (
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/4uvirik/ProductService/internal/entity"

	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

type Config struct {
	App      AppConfig      `yaml:"app"`
	Database DatabaseConfig `yaml:"database"`
	Logger   LoggerConfig   `yaml:"logger"`
}

type AppConfig struct {
	Name string `env:"APP_NAME" yaml:"name"`
	Host string `env:"APP_HOST" yaml:"host"`
	Port string `env:"APP_PORT" yaml:"port"`
}

type DatabaseConfig struct {
	Host     string `env:"DB_HOST"     yaml:"host"`
	Port     string `env:"DB_PORT"     yaml:"port"`
	User     string `env:"DB_USER"     yaml:"user"`
	Password string `env:"DB_PASSWORD" yaml:"password"`
	Name     string `env:"DB_NAME"     yaml:"name"`
	SSLMode  string `env:"DB_SSL_MODE" envDefault:"disable" yaml:"ssl_mode"`
}

type LoggerConfig struct {
	Level string `env:"LOG_LEVEL" envDefault:"info" yaml:"level"`
}

func loadFromEnv(config *Config) error {
	if _, err := os.Stat(".env"); err == nil {
		_ = godotenv.Load(".env")
	}

	if err := env.Parse(config); err != nil {
		return fmt.Errorf("error to parse env in config: %w", err)
	}

	return nil
}

func loadFromYaml(path string, config *Config) error {
	date, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read yaml file: %w", err)
	}

	if err := yaml.Unmarshal(date, config); err != nil {
		return fmt.Errorf("cant yaml unmarshal: %w", err)
	}

	return nil
}

// LoadConfig - загрузка конфига сначала из .env потом из yaml.
func LoadConfig(yamlPath string) (*Config, error) {
	config := &Config{}

	if err := loadFromEnv(config); err == nil {
		return config, nil
	}

	if err := loadFromYaml(yamlPath, config); err != nil {
		return nil, err
	}

	return config, nil
}

func (cfg *Config) DSN() string {
	hostAndPort := net.JoinHostPort(cfg.Database.Host, cfg.Database.Port)

	return fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		cfg.Database.User, cfg.Database.Password, hostAndPort, cfg.Database.Name,
	)
}

func (cfg *Config) Validate() error {
	var errorsMsg []string

	if cfg.App.Name == "" {
		errorsMsg = append(errorsMsg, "invalid app name")
	}

	if cfg.App.Host == "" {
		errorsMsg = append(errorsMsg, "invalid app host")
	}

	if cfg.App.Port == "" {
		errorsMsg = append(errorsMsg, "invalid app port")
	}

	if cfg.Database.Host == "" {
		errorsMsg = append(errorsMsg, "invalid database host")
	}

	if cfg.Database.Port == "" {
		errorsMsg = append(errorsMsg, "invalid database port")
	}

	if cfg.Database.User == "" {
		errorsMsg = append(errorsMsg, "invalid database user")
	}

	if cfg.Database.Password == "" {
		errorsMsg = append(errorsMsg, "invalid database password")
	}

	if cfg.Database.Name == "" {
		errorsMsg = append(errorsMsg, "invalid database name")
	}

	if len(errorsMsg) > 0 {
		return fmt.Errorf("%w: %s", entity.ErrInvalidConfig, strings.Join(errorsMsg, ": "))
	}

	return nil
}
