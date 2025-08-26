package config

import (
	"fmt"
	"os"
	"strings"

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
	Name string `yaml:"name" env:"APP_NAME"`
	Host string `yaml:"host" env:"APP_HOST"`
	Port string `yaml:"port" env:"APP_PORT"`
}

type DatabaseConfig struct {
	Host     string `yaml:"host" env:"DB_HOST"`
	Port     string `yaml:"port" env:"DB_PORT"`
	User     string `yaml:"user" env:"DB_USER"`
	Password string `yaml:"password" env:"DB_PASSWORD"`
	Name     string `yaml:"name" env:"DB_NAME"`
	SSLMode  string `yaml:"ssl_mode" env:"DB_SSL_MODE" envDefault:"disable"`
}

type LoggerConfig struct {
	Level string `yaml:"level" env:"LOG_LEVEL" envDefault:"info"`
}

func loadFromEnv(config *Config) error {

	if err := godotenv.Load(".env"); err != nil {
		return fmt.Errorf("cant load env file: %w", err)
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

func (c *Config) DSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		c.Database.User, c.Database.Password, c.Database.Host, c.Database.Port, c.Database.Name,
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
		return fmt.Errorf("%s", strings.Join(errorsMsg, ": "))
	}
	return nil
}
