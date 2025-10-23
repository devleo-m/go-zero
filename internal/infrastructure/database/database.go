package database

import (
	"fmt"
	"time"

	"github.com/devleo-m/go-zero/internal/infrastructure/config"
	"github.com/devleo-m/go-zero/internal/infrastructure/persistence/postgres"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// DB é a instância global do banco de dados
var DB *gorm.DB

// NewPostgresDB cria uma nova conexão com PostgreSQL
func NewPostgresDB(cfg *config.Config) *gorm.DB {
	// Usar a função de conexão que já criamos
	db, err := postgres.NewConnection()
	if err != nil {
		panic(fmt.Sprintf("failed to connect to database: %v", err))
	}

	// Configurar pool de conexões
	sqlDB, err := db.DB()
	if err != nil {
		panic(fmt.Sprintf("failed to get database instance: %v", err))
	}

	// Configurar pool de conexões
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db
}

// TestConnection testa a conexão com o banco
func TestConnection() error {
	if DB == nil {
		return fmt.Errorf("database not initialized")
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	return nil
}

// CloseDB fecha a conexão com o banco
func CloseDB() {
	if DB == nil {
		return
	}

	sqlDB, err := DB.DB()
	if err != nil {
		zap.L().Error("failed to get database instance for closing", zap.Error(err))
		return
	}

	if err := sqlDB.Close(); err != nil {
		zap.L().Error("failed to close database connection", zap.Error(err))
	}
}

// GetDB retorna a instância do banco de dados
func GetDB() *gorm.DB {
	return DB
}

// SetDB define a instância do banco de dados
func SetDB(db *gorm.DB) {
	DB = db
}
