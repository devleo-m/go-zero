package database

import (
	"fmt"
	"log"
	"time"

	"github.com/devleo-m/go-zero/internal/infra/config"
	"github.com/devleo-m/go-zero/internal/infra/logger"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// DB é a instância global do GORM
var DB *gorm.DB

// NewPostgresDB inicializa e retorna a conexão com o PostgreSQL
func NewPostgresDB(cfg *config.Config) *gorm.DB {
	dsn := cfg.DSN()

	// Configuração do GORM (opcional, mas recomendado)
	gormConfig := &gorm.Config{
		// Permite usar a nomenclatura Domain Driven Design (DDD) no Go
		// e snake_case no banco de dados (UserEntity -> user_entities)
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		// Configurações de logging do GORM (opcional)
		// Logger: logger.NewGormLogger(logger.Logger),
	}

	var err error

	// Tentativa de conexão com retries (essencial em containers, onde o DB pode subir depois da API)
	for i := 0; i < 5; i++ {
		DB, err = gorm.Open(postgres.Open(dsn), gormConfig)
		if err == nil {
			logger.Info("Conexão com o PostgreSQL estabelecida com sucesso!")
			// Configura o Pool de Conexões (Melhores Práticas de Performance)
			sqlDB, err := DB.DB()
			if err != nil {
				logger.Error("Erro ao obter o pool de conexões SQL", err)
				return nil
			}
			sqlDB.SetMaxIdleConns(10)           // Máximo de conexões ociosas
			sqlDB.SetMaxOpenConns(100)          // Máximo de conexões abertas
			sqlDB.SetConnMaxLifetime(time.Hour) // Tempo máximo de vida de uma conexão
			return DB
		}

		logger.Error(fmt.Sprintf("Falha ao conectar ao PostgreSQL, tentando novamente em 5 segundos... (Tentativa %d/5)", i+1), err, zap.String("dsn", dsn))
		time.Sleep(5 * time.Second)
	}

	log.Fatal("Não foi possível conectar ao banco de dados após várias tentativas.")
	return nil // Nunca alcançado devido ao Fatal
}

// TestConnection testa a conexão com o banco de dados
func TestConnection() error {
	if DB == nil {
		return fmt.Errorf("banco de dados não inicializado")
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("erro ao obter conexão SQL: %w", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("falha ao fazer ping no banco: %w", err)
	}

	return nil
}

// CloseDB fecha a conexão com o banco de dados (útil para graceful shutdown)
func CloseDB() {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err == nil {
			sqlDB.Close()
			logger.Info("Conexão com o PostgreSQL fechada.")
		}
	}
}
