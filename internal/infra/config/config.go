package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Config representa todas as configurações da aplicação
type Config struct {
	// Aplicação
	App AppConfig `mapstructure:"app"`

	// Banco de Dados
	Database DatabaseConfig `mapstructure:"database"`

	// Redis
	Redis RedisConfig `mapstructure:"redis"`

	// Logs
	Logger LoggerConfig `mapstructure:"logger"`

	// JWT
	JWT JWTConfig `mapstructure:"jwt"`

	// Serviços Externos
	External ExternalConfig `mapstructure:"external"`
}

// AppConfig configurações da aplicação
type AppConfig struct {
	Name string `mapstructure:"name"`
	Env  string `mapstructure:"env"`
	Port string `mapstructure:"port"`
}

// DatabaseConfig configurações do banco de dados
type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Name     string `mapstructure:"name"`
	SSLMode  string `mapstructure:"ssl_mode"`
}

// RedisConfig configurações do Redis
type RedisConfig struct {
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
}

// LoggerConfig configurações do logger
type LoggerConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
}

// JWTConfig configurações do JWT
type JWTConfig struct {
	Secret           string        `mapstructure:"secret"`
	ExpiresIn        time.Duration `mapstructure:"expires_in"`
	RefreshExpiresIn time.Duration `mapstructure:"refresh_expires_in"`
}

// ExternalConfig configurações de serviços externos
type ExternalConfig struct {
	Stripe   StripeConfig   `mapstructure:"stripe"`
	SendGrid SendGridConfig `mapstructure:"sendgrid"`
	MinIO    MinIOConfig    `mapstructure:"minio"`
}

type StripeConfig struct {
	SecretKey string `mapstructure:"secret_key"`
}

type SendGridConfig struct {
	APIKey string `mapstructure:"api_key"`
}

type MinIOConfig struct {
	Endpoint  string `mapstructure:"endpoint"`
	AccessKey string `mapstructure:"access_key"`
	SecretKey string `mapstructure:"secret_key"`
}

// LoadConfig carrega e valida as configurações
func LoadConfig() (*Config, error) {
	// 1. Configurar o Viper
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")

	// 2. Configurar valores padrão
	setDefaults()

	// 3. Permitir variáveis de ambiente
	viper.AutomaticEnv()

	// 4. Ler arquivo de configuração (opcional)
	if err := viper.ReadInConfig(); err != nil {
		// Não é fatal se não encontrar o arquivo
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("erro ao ler arquivo de configuração: %w", err)
		}
	}

	// 5. Desserializar para struct
	var config Config

	// Mapear variáveis de ambiente para struct aninhada
	config.App.Name = viper.GetString("APP_NAME")
	config.App.Env = viper.GetString("APP_ENV")
	config.App.Port = viper.GetString("APP_PORT")

	config.Database.Host = viper.GetString("DB_HOST")
	config.Database.Port = viper.GetString("DB_PORT")
	config.Database.User = viper.GetString("DB_USER")
	config.Database.Password = viper.GetString("DB_PASSWORD")
	config.Database.Name = viper.GetString("DB_NAME")
	config.Database.SSLMode = viper.GetString("DB_SSLMODE")

	config.Redis.Host = viper.GetString("REDIS_HOST")
	config.Redis.Port = viper.GetString("REDIS_PORT")

	config.Logger.Level = viper.GetString("LOG_LEVEL")
	config.Logger.Format = viper.GetString("LOG_FORMAT")

	config.JWT.Secret = viper.GetString("JWT_SECRET")
	config.JWT.ExpiresIn, _ = time.ParseDuration(viper.GetString("JWT_EXPIRES_IN"))
	config.JWT.RefreshExpiresIn, _ = time.ParseDuration(viper.GetString("REFRESH_TOKEN_EXPIRES_IN"))

	config.External.Stripe.SecretKey = viper.GetString("STRIPE_SECRET_KEY")
	config.External.SendGrid.APIKey = viper.GetString("SENDGRID_API_KEY")
	config.External.MinIO.Endpoint = viper.GetString("MINIO_ENDPOINT")
	config.External.MinIO.AccessKey = viper.GetString("MINIO_ACCESS_KEY")
	config.External.MinIO.SecretKey = viper.GetString("MINIO_SECRET_KEY")

	// 6. Validar configurações
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("configurações inválidas: %w", err)
	}

	return &config, nil
}

// setDefaults define valores padrão para todas as configurações
func setDefaults() {
	// Aplicação
	viper.SetDefault("APP_NAME", "go-zero")
	viper.SetDefault("APP_ENV", "development")
	viper.SetDefault("APP_PORT", "8080")

	// Banco de Dados
	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_PORT", "5432")
	viper.SetDefault("DB_USER", "postgres")
	viper.SetDefault("DB_PASSWORD", "postgres123")
	viper.SetDefault("DB_NAME", "go_zero_dev")
	viper.SetDefault("DB_SSLMODE", "disable")

	// Redis
	viper.SetDefault("REDIS_HOST", "localhost")
	viper.SetDefault("REDIS_PORT", "6379")

	// Logger
	viper.SetDefault("LOG_LEVEL", "debug")
	viper.SetDefault("LOG_FORMAT", "json")

	// JWT
	viper.SetDefault("JWT_SECRET", "your-super-secret-jwt-key-change-in-production")
	viper.SetDefault("JWT_EXPIRES_IN", "24h")
	viper.SetDefault("REFRESH_TOKEN_EXPIRES_IN", "168h")

	// Serviços Externos
	viper.SetDefault("STRIPE_SECRET_KEY", "")
	viper.SetDefault("SENDGRID_API_KEY", "")
	viper.SetDefault("MINIO_ENDPOINT", "localhost:9000")
	viper.SetDefault("MINIO_ACCESS_KEY", "minioadmin")
	viper.SetDefault("MINIO_SECRET_KEY", "minioadmin123")
}

// Validate valida se todas as configurações obrigatórias estão presentes
func (c *Config) Validate() error {
	var errors []string

	// Validar aplicação
	if c.App.Name == "" {
		errors = append(errors, "APP_NAME é obrigatório")
	}
	if c.App.Env == "" {
		errors = append(errors, "APP_ENV é obrigatório")
	}
	if c.App.Port == "" {
		errors = append(errors, "APP_PORT é obrigatório")
	}

	// Validar banco de dados
	if c.Database.Host == "" {
		errors = append(errors, "DB_HOST é obrigatório")
	}
	if c.Database.User == "" {
		errors = append(errors, "DB_USER é obrigatório")
	}
	if c.Database.Password == "" {
		errors = append(errors, "DB_PASSWORD é obrigatório")
	}
	if c.Database.Name == "" {
		errors = append(errors, "DB_NAME é obrigatório")
	}

	// Validar JWT
	if c.JWT.Secret == "" {
		errors = append(errors, "JWT_SECRET é obrigatório")
	}
	if c.JWT.Secret == "your-super-secret-jwt-key-change-in-production" && c.App.Env == "production" {
		errors = append(errors, "JWT_SECRET deve ser alterado em produção")
	}

	// Validar Redis
	if c.Redis.Host == "" {
		errors = append(errors, "REDIS_HOST é obrigatório")
	}

	if len(errors) > 0 {
		return fmt.Errorf("configurações inválidas: %s", strings.Join(errors, ", "))
	}

	return nil
}

// DSN retorna a string de conexão do PostgreSQL
func (c *Config) DSN() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=America/Sao_Paulo",
		c.Database.Host,
		c.Database.User,
		c.Database.Password,
		c.Database.Name,
		c.Database.Port,
		c.Database.SSLMode,
	)
}

// RedisURL retorna a URL de conexão do Redis
func (c *Config) RedisURL() string {
	return fmt.Sprintf("%s:%s", c.Redis.Host, c.Redis.Port)
}

// IsDevelopment verifica se está em ambiente de desenvolvimento
func (c *Config) IsDevelopment() bool {
	return c.App.Env == "development"
}

// IsProduction verifica se está em ambiente de produção
func (c *Config) IsProduction() bool {
	return c.App.Env == "production"
}
