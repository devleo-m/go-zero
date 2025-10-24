package main

import (
	"log"
	"os"

	"github.com/devleo-m/go-zero/internal/infrastructure"
	userApp "github.com/devleo-m/go-zero/internal/modules/user/application"
	userRepo "github.com/devleo-m/go-zero/internal/modules/user/infrastructure/postgres"
	userHttp "github.com/devleo-m/go-zero/internal/modules/user/presentation/http"
	"github.com/gin-gonic/gin"
)

func main() {
	// Configurar Gin
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	// Conectar ao banco de dados
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "host=localhost user=postgres password=postgres dbname=go_zero port=5432 sslmode=disable"
	}

	db, err := infrastructure.NewDatabase(dsn)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

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

	// Configurar rotas
	userHttp.SetupRoutes(router, userHandler)

	// Rota de health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Iniciar servidor
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
