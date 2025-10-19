package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	// Aplicação
	AppName string `mapstructure:"APP_NAME"`
	AppEnv  string `mapstructure:"APP_ENV"`
	AppPort string `mapstructure:"APP_PORT"`

	// Banco de Dados
	DBHost    string `mapstructure:"DB_HOST"`
	DBPort    string `mapstructure:"DB_PORT"`
	DBUser    string `mapstructure:"DB_USER"`
	DBPass    string `mapstructure:"DB_PASSWORD"`
	DBName    string `mapstructure:"DB_NAME"`
	DBSslMode string `mapstructure:"DB_SSLMODE"`

	// Redis
	RedisHost string `mapstructure:"REDIS_HOST"`
	RedisPort string `mapstructure:"REDIS_PORT"`

	// Logs
	LogLevel  string `mapstructure:"LOG_LEVEL"`
	LogFormat string `mapstructure:"LOG_FORMAT"`

	// JWT
	JWTSecret           string `mapstructure:"JWT_SECRET"`
	JWTExpiresIn        string `mapstructure:"JWT_EXPIRES_IN"`
	RefreshTokenExpires string `mapstructure:"REFRESH_TOKEN_EXPIRES_IN"`
}

// LoadConfig carrega as configurações do arquivo .env ou variáveis de ambiente
func LoadConfig() (config Config) {
	// 1. Configurar o Viper para ler de .env (se existir)
	viper.AddConfigPath(".") // Procura na raiz
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	// 2. Configurar valores padrão
	viper.SetDefault("APP_NAME", "go-zero")
	viper.SetDefault("APP_ENV", "development")
	viper.SetDefault("APP_PORT", "8080")
	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_PORT", "5432")
	viper.SetDefault("DB_USER", "postgres")
	viper.SetDefault("DB_PASSWORD", "postgres123")
	viper.SetDefault("DB_NAME", "go_zero_dev")
	viper.SetDefault("DB_SSLMODE", "disable")
	viper.SetDefault("REDIS_HOST", "localhost")
	viper.SetDefault("REDIS_PORT", "6379")
	viper.SetDefault("LOG_LEVEL", "debug")
	viper.SetDefault("LOG_FORMAT", "json")
	viper.SetDefault("JWT_SECRET", "your-super-secret-jwt-key-change-in-production")
	viper.SetDefault("JWT_EXPIRES_IN", "24h")
	viper.SetDefault("REFRESH_TOKEN_EXPIRES_IN", "168h")

	// 3. Sobrescrever com variáveis de ambiente (útil no Docker)
	viper.AutomaticEnv()

	// 4. Tenta ler o arquivo .env
	if err := viper.ReadInConfig(); err != nil {
		// Não é um erro fatal se o .env não for encontrado, pois podemos usar apenas env vars
		log.Println("WARNING: Could not load .env file. Using environment variables and defaults.")
	}

	// 5. Desserializar para a struct Config usando mapstructure
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("Fatal error config file: %s \n", err)
	}

	return
}

// DSN retorna a string de conexão do PostgreSQL
func (c *Config) DSN() string {
	return "host=" + c.DBHost +
		" user=" + c.DBUser +
		" password=" + c.DBPass +
		" dbname=" + c.DBName +
		" port=" + c.DBPort +
		" sslmode=" + c.DBSslMode +
		" TimeZone=America/Sao_Paulo" // Exemplo: Timezone
}
