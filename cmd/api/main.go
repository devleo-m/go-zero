package main

import (
	"log"

	"github.com/devleo-m/go-zero/internal/infrastructure"
	"github.com/devleo-m/go-zero/internal/infrastructure/config"
	"github.com/devleo-m/go-zero/internal/infrastructure/http/middleware"
	"github.com/devleo-m/go-zero/internal/infrastructure/http/routes"
	"github.com/devleo-m/go-zero/internal/infrastructure/logger"
	userApp "github.com/devleo-m/go-zero/internal/modules/user/application"
	userRepo "github.com/devleo-m/go-zero/internal/modules/user/infrastructure/postgres"
	userHttp "github.com/devleo-m/go-zero/internal/modules/user/presentation/http"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	// Carregar configurações
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Configurar logger
	appLogger, err := logger.NewFromEnv()
	if err != nil {
		log.Fatal("Failed to initialize logger:", err)
	}
	defer appLogger.Sync()

	appLogger.Info("Starting GO ZERO API",
		zap.String("app_name", cfg.App.Name),
		zap.String("version", cfg.App.Version),
		zap.String("environment", cfg.App.Env),
	)

	// Configurar Gin
	if cfg.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	// Conectar ao banco de dados
	dsn := cfg.Database.URL
	if dsn == "" {
		dsn = buildDatabaseURL(cfg.Database)
	}

	db, err := infrastructure.NewDatabase(dsn)
	if err != nil {
		appLogger.Fatal("Failed to connect to database",
			zap.Error(err),
			zap.String("component", "database"),
		)
	}
	defer db.Close()

	appLogger.Info("Database connected successfully",
		zap.String("component", "database"),
	)

	// Configurar repositórios
	userRepository := userRepo.NewRepository(db.DB)

	// Configurar use cases
	createUserUseCase := userApp.NewCreateUserUseCase(userRepository)
	getUserUseCase := userApp.NewGetUserUseCase(userRepository)
	listUsersUseCase := userApp.NewListUsersUseCase(userRepository)
	updateUserUseCase := userApp.NewUpdateUserUseCase(userRepository)
	deleteUserUseCase := userApp.NewDeleteUserUseCase(userRepository)

	// Configurar handlers
	userHandler := userHttp.NewHandler(
		createUserUseCase,
		getUserUseCase,
		listUsersUseCase,
		updateUserUseCase,
		deleteUserUseCase,
	)

	// Configurar rate limiter
	rateLimiter := middleware.NewRateLimiter(cfg.RateLimit.Requests, cfg.RateLimit.Window)

	// Configurar rotas
	router := gin.New()
	routesConfig := &routes.Config{
		JWT: routes.JWTConfig{
			Secret: cfg.JWT.Secret,
		},
		CORS: routes.CORSConfig{
			AllowedOrigins: cfg.CORS.AllowedOrigins,
			AllowedMethods: cfg.CORS.AllowedMethods,
			AllowedHeaders: cfg.CORS.AllowedHeaders,
		},
		RateLimiter: rateLimiter,
		UserHandler: userHandler,
	}

	routes.SetupRoutes(router, routesConfig)

	// Iniciar servidor
	port := cfg.App.Port
	if port == "" {
		port = "8080"
	}

	appLogger.Info("Server starting",
		zap.String("port", port),
		zap.String("component", "server"),
	)

	if err := router.Run(":" + port); err != nil {
		appLogger.Fatal("Failed to start server",
			zap.Error(err),
			zap.String("component", "server"),
		)
	}
}

// buildDatabaseURL constrói a URL do banco de dados
func buildDatabaseURL(dbCfg config.DatabaseConfig) string {
	return "host=" + dbCfg.Host +
		" user=" + dbCfg.User +
		" password=" + dbCfg.Password +
		" dbname=" + dbCfg.Name +
		" port=" + dbCfg.Port +
		" sslmode=" + dbCfg.SSLMode
}
