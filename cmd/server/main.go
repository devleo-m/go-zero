package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"

	"go.uber.org/zap"

	"github.com/devleo-m/go-zero/internal/infrastructure"
	"github.com/devleo-m/go-zero/internal/infrastructure/logger"
	userApp "github.com/devleo-m/go-zero/internal/modules/user/application"
	userHttp "github.com/devleo-m/go-zero/internal/modules/user/infrastructure/http"
	userRepo "github.com/devleo-m/go-zero/internal/modules/user/infrastructure/postgres"
)

const (
	defaultDSN  = "host=localhost user=postgres password=postgres dbname=go_zero port=5432 sslmode=disable"
	defaultPort = "8080"
)

func main() {
	// Configurar logger
	appLogger, err := logger.NewFromEnv()
	if err != nil {
		log.Fatal("Failed to initialize logger:", err)
	}

	defer func() {
		if syncErr := appLogger.Sync(); syncErr != nil {
			log.Printf("Failed to sync logger: %v", syncErr)
		}
	}()

	// Configurar Gin
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	// Conectar ao banco de dados
	dsn := getDatabaseDSN()

	db, err := infrastructure.NewDatabase(dsn)
	if err != nil {
		appLogger.Fatal("Failed to connect to database", zap.Error(err))
	}

	defer func() {
		if closeErr := db.Close(); closeErr != nil {
			appLogger.Error("Failed to close database connection", zap.Error(closeErr))
		}
	}()

	// Configurar módulo de usuários
	setupUserModule(router, db)

	// Configurar rota de health check
	setupHealthCheck(router)

	// Iniciar servidor
	port := getServerPort()
	appLogger.Info("Server starting", zap.String("port", port))

	if err := router.Run(":" + port); err != nil {
		appLogger.Fatal("Failed to start server", zap.Error(err))
	}
}

// getDatabaseDSN obtém a DSN do banco de dados.
func getDatabaseDSN() string {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = defaultDSN
	}

	return dsn
}

// getServerPort obtém a porta do servidor.
func getServerPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	return port
}

// setupUserModule configura o módulo de usuários.
func setupUserModule(router *gin.Engine, db *infrastructure.Database) {
	userRepository := userRepo.NewRepository(db.DB)

	createUserUseCase := userApp.NewCreateUserUseCase(userRepository)
	getUserUseCase := userApp.NewGetUserUseCase(userRepository)
	listUsersUseCase := userApp.NewListUsersUseCase(userRepository)
	updateUserUseCase := userApp.NewUpdateUserUseCase(userRepository)
	deleteUserUseCase := userApp.NewDeleteUserUseCase(userRepository)

	userHandler := userHttp.NewHandler(
		createUserUseCase,
		getUserUseCase,
		listUsersUseCase,
		updateUserUseCase,
		deleteUserUseCase,
	)

	userHttp.SetupRoutes(router, userHandler)
}

// setupHealthCheck configura a rota de health check.
func setupHealthCheck(router *gin.Engine) {
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
}
