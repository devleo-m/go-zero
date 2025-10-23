package config

import (
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

// Config representa todas as configurações da aplicação
type Config struct {
	App      AppConfig
	Database DatabaseConfig
	Logger   LoggerConfig
}

// AppConfig configurações da aplicação
type AppConfig struct {
	Name    string
	Env     string
	Port    string
	Version string
}

// DatabaseConfig configurações do banco de dados
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

// LoggerConfig configurações do logger
type LoggerConfig struct {
	Level  string
	Format string
}

// LoadConfig carrega configurações do ambiente
func LoadConfig() (*Config, error) {
	// Carregar .env se existir
	_ = godotenv.Load()

	return &Config{
		App: AppConfig{
			Name:    getEnv("APP_NAME", "GO ZERO"),
			Env:     getEnv("APP_ENV", "development"),
			Port:    getEnv("APP_PORT", "8080"),
			Version: getEnv("APP_VERSION", "1.0.0"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			Name:     getEnv("DB_NAME", "go_zero"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		Logger: LoggerConfig{
			Level:  getEnv("LOG_LEVEL", "info"),
			Format: getEnv("LOG_FORMAT", "json"),
		},
	}, nil
}

// IsProduction verifica se está em produção
func (c *Config) IsProduction() bool {
	return strings.ToLower(c.App.Env) == "production"
}

// IsDevelopment verifica se está em desenvolvimento
func (c *Config) IsDevelopment() bool {
	return strings.ToLower(c.App.Env) == "development"
}

// IsTest verifica se está em teste
func (c *Config) IsTest() bool {
	return strings.ToLower(c.App.Env) == "test"
}

// getEnv obtém variável de ambiente com valor padrão
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt obtém variável de ambiente como int
func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// getEnvAsBool obtém variável de ambiente como bool
func getEnvAsBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}
