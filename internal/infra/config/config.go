package config

import (
	"log"
	"os"
	"github.com/spf13/viper"
)

type Config struct {
	AppPort string
	DBHost string
	DBPort string
	DBUser string
	DBPass string
	DBName string
	DBSslMode string
}

// LoadConfig carrega as configurações do arquivo .env ou variáveis de ambiente
func LoadConfig() (config Config) {
	// 1. Configurar o Viper para ler de .env (se existir)
	viper.AddConfigPath(".") // Procura na raiz
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	// 2. Sobrescrever com variáveis de ambiente (útil no Docker)
	viper.AutomaticEnv()

	// Tenta ler o arquivo .env
	if err := viper.ReadInConfig(); err != nil {
		// Não é um erro fatal se o .env não for encontrado, pois podemos usar apenas env vars
		log.Println("WARNING: Could not load .env file. Using environment variables only.")
	}

	// 3. Desserializar para a struct Config
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("Fatal error config file: %s \n", err)
	}
	
	// Preenche campos não encontrados com fallback de ENV para garantir
	if config.DBHost == "" {
        config.DBHost = os.Getenv("DB_HOST")
    }
    // ... repetir para outros campos se necessário, mas viper.AutomaticEnv() já deve cuidar disso

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