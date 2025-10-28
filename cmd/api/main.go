package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/devleo-m/go-zero/internal/infrastructure"
	"github.com/devleo-m/go-zero/internal/infrastructure/config"
	"github.com/devleo-m/go-zero/internal/infrastructure/http/middleware"
	"github.com/devleo-m/go-zero/internal/infrastructure/http/routes"
	"github.com/devleo-m/go-zero/internal/infrastructure/logger"
	userApp "github.com/devleo-m/go-zero/internal/modules/user/application"
	userHttp "github.com/devleo-m/go-zero/internal/modules/user/infrastructure/http"
	userRepo "github.com/devleo-m/go-zero/internal/modules/user/infrastructure/postgres"
)

func main() {
	// Carregar configurações
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Configurar logger
	appLogger := setupLogger()
	defer syncLogger(appLogger)

	logAppStart(appLogger, cfg)

	// Configurar Gin mode
	configureGinMode(cfg.App.Env)

	// Conectar ao banco de dados
	db := setupDatabase(cfg, appLogger)
	defer closeDatabase(db, appLogger)

	appLogger.Info("Database connected successfully",
		zap.String("component", "database"),
	)

	// Configurar handlers e rotas
	router := setupRouter(cfg, db)

	// Iniciar servidor
	startServer(router, cfg.App.Port, appLogger)
}

// setupLogger inicializa e retorna o logger.
func setupLogger() *logger.Logger {
	appLogger, err := logger.NewFromEnv()
	if err != nil {
		log.Fatal("Failed to initialize logger:", err)
	}

	return appLogger
}

// syncLogger sincroniza o logger antes de encerrar.
func syncLogger(appLogger *logger.Logger) {
	if syncErr := appLogger.Sync(); syncErr != nil {
		log.Printf("Failed to sync logger: %v", syncErr)
	}
}

// logAppStart registra o início da aplicação.
func logAppStart(appLogger *logger.Logger, cfg *config.Config) {
	appLogger.Info("Starting GO ZERO API",
		zap.String("app_name", cfg.App.Name),
		zap.String("version", cfg.App.Version),
		zap.String("environment", cfg.App.Env),
	)
}

// configureGinMode configura o modo do Gin baseado no ambiente.
func configureGinMode(env string) {
	if env == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}
}

// setupDatabase conecta ao banco de dados e retorna a conexão.
func setupDatabase(cfg *config.Config, appLogger *logger.Logger) *infrastructure.Database {
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

	return db
}

// closeDatabase fecha a conexão com o banco de dados.
func closeDatabase(db *infrastructure.Database, appLogger *logger.Logger) {
	if closeErr := db.Close(); closeErr != nil {
		appLogger.Error("Failed to close database connection",
			zap.Error(closeErr),
			zap.String("component", "database"),
		)
	}
}

// setupRouter configura e retorna o router com todas as rotas.
func setupRouter(cfg *config.Config, db *infrastructure.Database) *gin.Engine {
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

	return router
}

// startServer inicia o servidor HTTP.
func startServer(router *gin.Engine, port string, appLogger *logger.Logger) {
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

// buildDatabaseURL constrói a URL do banco de dados.
func buildDatabaseURL(dbCfg config.DatabaseConfig) string {
	return "host=" + dbCfg.Host +
		" user=" + dbCfg.User +
		" password=" + dbCfg.Password +
		" dbname=" + dbCfg.Name +
		" port=" + dbCfg.Port +
		" sslmode=" + dbCfg.SSLMode
}
