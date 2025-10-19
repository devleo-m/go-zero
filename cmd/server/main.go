package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/devleo-m/go-zero/internal/infra/config"
	"github.com/devleo-m/go-zero/internal/infra/database"
	"github.com/devleo-m/go-zero/internal/infra/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	// 1. Carregar e validar configurações
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("❌ Erro ao carregar configurações: %v\n", err)
		os.Exit(1)
	}

	// 2. Inicializar logger
	if err := logger.InitLogger(logger.Config{
		Level:  cfg.Logger.Level,
		Format: cfg.Logger.Format,
	}); err != nil {
		fmt.Printf("❌ Erro ao inicializar logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	// 3. Log de inicialização
	logger.Info("🚀 Iniciando GO ZERO API",
		zap.String("app_name", cfg.App.Name),
		zap.String("app_env", cfg.App.Env),
		zap.String("app_port", cfg.App.Port),
		zap.String("log_level", cfg.Logger.Level),
	)

	// 4. Inicializar banco de dados
	if err := initDatabase(cfg); err != nil {
		logger.Fatal("❌ Falha ao conectar com banco de dados", err)
	}
	defer database.CloseDB()

	// 5. Inicializar servidor HTTP
	server := initHTTPServer(cfg)

	// 6. Configurar graceful shutdown
	gracefulShutdown(server, cfg)
}

// initDatabase inicializa conexão com banco de dados
func initDatabase(cfg *config.Config) error {
	logger.Info("🔌 Conectando com banco de dados...",
		zap.String("host", cfg.Database.Host),
		zap.String("port", cfg.Database.Port),
		zap.String("database", cfg.Database.Name),
	)

	// Inicializar banco
	database.DB = database.NewPostgresDB(cfg)

	// Testar conexão
	if err := database.TestConnection(); err != nil {
		return fmt.Errorf("falha ao testar conexão: %w", err)
	}

	logger.Info("✅ Banco de dados conectado com sucesso!")
	return nil
}

// initHTTPServer inicializa servidor HTTP
func initHTTPServer(cfg *config.Config) *http.Server {
	// Configurar Gin baseado no ambiente
	if cfg.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	// Criar router
	router := gin.New()

	// Middlewares globais
	router.Use(gin.Recovery())
	router.Use(loggerMiddleware())

	// Rotas
	setupRoutes(router, cfg)

	// Criar servidor HTTP
	server := &http.Server{
		Addr:    ":" + cfg.App.Port,
		Handler: router,
	}

	// Iniciar servidor em goroutine
	go func() {
		logger.Info("🌐 Servidor HTTP iniciado",
			zap.String("port", cfg.App.Port),
			zap.String("env", cfg.App.Env),
		)

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("❌ Erro ao iniciar servidor HTTP", err)
		}
	}()

	return server
}

// setupRoutes configura todas as rotas da aplicação
func setupRoutes(router *gin.Engine, cfg *config.Config) {
	// Health check
	router.GET("/health", healthCheckHandler(cfg))

	// API v1
	v1 := router.Group("/api/v1")
	{
		// Rotas públicas
		v1.GET("/status", statusHandler)

		// Rotas protegidas (serão implementadas futuramente)
		// v1.Use(authMiddleware())
		// v1.GET("/users", usersHandler)
	}

	logger.Info("📋 Rotas configuradas com sucesso")
}

// healthCheckHandler handler para health check
func healthCheckHandler(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Verificar status do banco
		dbStatus := "healthy"
		if err := database.TestConnection(); err != nil {
			dbStatus = "unhealthy"
			logger.Error("Health check: banco indisponível", err)
		}

		// Preparar resposta
		response := gin.H{
			"status":    "up",
			"service":   cfg.App.Name,
			"version":   "1.0.0",
			"env":       cfg.App.Env,
			"timestamp": time.Now().Format(time.RFC3339),
			"uptime":    time.Since(start).String(),
			"database":  dbStatus,
		}

		// Status code baseado na saúde do sistema
		statusCode := http.StatusOK
		if dbStatus == "unhealthy" {
			statusCode = http.StatusServiceUnavailable
		}

		// Log da requisição
		logger.LogHTTPRequest(
			c.Request.Method,
			c.Request.URL.Path,
			statusCode,
			time.Since(start),
			"", // userID (não disponível em health check)
		)

		c.JSON(statusCode, response)
	}
}

// statusHandler handler para status da API
func statusHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "GO ZERO API está funcionando! 🚀",
		"features": []string{
			"✅ Configurações validadas",
			"✅ Logger estruturado",
			"✅ Banco de dados conectado",
			"✅ Hot reload funcionando",
			"✅ Graceful shutdown",
		},
		"next_steps": []string{
			"🔐 Implementar autenticação JWT",
			"👥 Criar módulo de usuários",
			"🛒 Implementar e-commerce",
			"📚 Adicionar sistema de cursos",
		},
	})
}

// loggerMiddleware middleware para logging de requisições
func loggerMiddleware() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// Usar nosso logger estruturado
		logger.LogHTTPRequest(
			param.Method,
			param.Path,
			param.StatusCode,
			param.Latency,
			"", // userID (será implementado com auth)
		)
		return ""
	})
}

// gracefulShutdown configura shutdown graceful
func gracefulShutdown(server *http.Server, cfg *config.Config) {
	// Canal para receber sinais do sistema
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Aguardar sinal de shutdown
	<-quit
	logger.Info("🛑 Iniciando shutdown graceful...")

	// Contexto com timeout para shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Tentar shutdown graceful
	if err := server.Shutdown(ctx); err != nil {
		logger.Error("❌ Erro durante shutdown graceful", err)
		os.Exit(1)
	}

	logger.Info("✅ Shutdown graceful concluído com sucesso!")
}
