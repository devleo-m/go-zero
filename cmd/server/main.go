package main

import (
	"fmt"
	"os"
	"time"

	"github.com/devleo-m/go-zero/internal/infra/config"
	"github.com/devleo-m/go-zero/internal/infra/database"
	"github.com/devleo-m/go-zero/internal/infra/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	// 1. Carregar Configurações
	cfg := config.LoadConfig()

	// 2. Inicializar Logger (usa APP_ENV, que deve ser setada)
	// Vamos setar um default temporário se não estiver no ENV
	appEnv := os.Getenv("APP_ENV")
	if appEnv == "" {
		appEnv = "development"
	}
	logger.InitLogger(appEnv)
	logger.Info("Configurações carregadas e Logger inicializado.", zap.String("app_env", appEnv))

	// 3. Inicializar Banco de Dados (Será implementado no próximo passo/fase)
	database.DB = database.NewPostgresDB(cfg)
	// Se for rodar local sem Docker, pode pular a inicialização do DB por enquanto.
	defer database.CloseDB()

	// 4. Inicializar o Gin
	router := gin.Default()

	// 5. Setup de Rotas
	router.GET("/health", func(c *gin.Context) {
		logger.Info("Health check realizado com sucesso.")
		c.JSON(200, gin.H{
			"status":    "up",
			"service":   "Go Hexagonal API - HOT RELOAD FUNCIONANDO! 🚀",
			"message":   "Mudança aplicada automaticamente pelo Air!",
			"timestamp": time.Now().Format("2006-01-02 15:04:05"),
		})
	})

	// 6. Rodar o Servidor
	logger.Info(fmt.Sprintf("Servidor rodando na porta: %s", cfg.AppPort))
	if err := router.Run(fmt.Sprintf(":%s", cfg.AppPort)); err != nil {
		logger.Fatal("Erro ao iniciar o servidor Gin", err)
	}
}
