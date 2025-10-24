package routes

import (
	"github.com/devleo-m/go-zero/internal/infrastructure/http/middleware"
	"github.com/devleo-m/go-zero/internal/shared/response"
	"github.com/gin-gonic/gin"
)

// SetupRoutes configura todas as rotas da aplicação
func SetupRoutes(router *gin.Engine, config *Config) {
	// Middleware global
	router.Use(middleware.LoggingMiddleware(nil))
	router.Use(middleware.RequestIDMiddleware())
	router.Use(middleware.RecoveryMiddleware())
	router.Use(middleware.CORS(middleware.CORSConfig{
		AllowedOrigins:   config.CORS.AllowedOrigins,
		AllowedMethods:   config.CORS.AllowedMethods,
		AllowedHeaders:   config.CORS.AllowedHeaders,
		MaxAge:           3600,
		AllowCredentials: true,
	}))

	// Rate limiting
	if config.RateLimiter != nil {
		if rateLimiter, ok := config.RateLimiter.(*middleware.RateLimiter); ok {
			router.Use(middleware.RateLimit(rateLimiter))
		}
	}

	// Health check
	router.GET("/health", healthCheck)
	router.GET("/metrics", metricsHandler)

	// API v1
	v1 := router.Group("/api/v1")
	{
		// Rotas públicas
		// public := v1.Group("/")
		// {
		// 	// Auth routes (será implementado)
		// 	// public.POST("/auth/register", authHandler.Register)
		// 	// public.POST("/auth/login", authHandler.Login)
		// }

		// Rotas protegidas
		protected := v1.Group("/")
		protected.Use(middleware.AuthMiddleware(config.JWT.Secret))
		{
			// User routes
			if config.UserHandler != nil {
				if userHandler, ok := config.UserHandler.(interface {
					CreateUser(*gin.Context)
					ListUsers(*gin.Context)
					GetUser(*gin.Context)
					UpdateUser(*gin.Context)
					DeleteUser(*gin.Context)
				}); ok {
					userRoutes := protected.Group("/users")
					{
						userRoutes.POST("", userHandler.CreateUser)
						userRoutes.GET("", userHandler.ListUsers)
						userRoutes.GET("/:id", userHandler.GetUser)
						userRoutes.PUT("/:id", userHandler.UpdateUser)
						userRoutes.DELETE("/:id", userHandler.DeleteUser)
					}
				}
			}

			// Admin routes
			admin := protected.Group("/admin")
			admin.Use(middleware.RequireRole("admin"))
			{
				// Admin-specific routes
				admin.GET("/stats", adminStats)
			}
		}
	}

	// WebSocket routes (será implementado)
	ws := router.Group("/ws")
	{
		ws.GET("/chat", chatWebSocket)
	}

	// Swagger documentation
	router.GET("/swagger/*any", swaggerHandler)
}

// healthCheck retorna o status de saúde da aplicação
func healthCheck(c *gin.Context) {
	response.Success(c, gin.H{
		"status":    "ok",
		"timestamp": gin.H{},
		"services": gin.H{
			"database": "ok",
			"redis":    "ok",
			"api":      "ok",
		},
	}, "Service is healthy")
}

// metricsHandler retorna métricas da aplicação
func metricsHandler(c *gin.Context) {
	// Aqui você pode implementar métricas customizadas
	// Por enquanto, retornamos um placeholder
	c.JSON(200, gin.H{
		"message": "Metrics endpoint - implement Prometheus metrics here",
	})
}

// adminStats retorna estatísticas administrativas
func adminStats(c *gin.Context) {
	response.Success(c, gin.H{
		"users": gin.H{
			"total":  0,
			"active": 0,
		},
		"system": gin.H{
			"uptime":  "0s",
			"version": "1.0.0",
		},
	}, "Admin statistics")
}

// chatWebSocket lida com conexões WebSocket para chat
func chatWebSocket(c *gin.Context) {
	// Implementação do WebSocket será feita posteriormente
	c.JSON(501, gin.H{
		"error": "WebSocket not implemented yet",
	})
}

// swaggerHandler serve a documentação Swagger
func swaggerHandler(c *gin.Context) {
	// Implementação do Swagger será feita posteriormente
	c.JSON(501, gin.H{
		"error": "Swagger documentation not implemented yet",
	})
}

// Config representa a configuração das rotas
type Config struct {
	JWT         JWTConfig
	CORS        CORSConfig
	RateLimiter interface{} // *middleware.RateLimiter
	UserHandler interface{} // *userHttp.Handler
}

type JWTConfig struct {
	Secret string
}

type CORSConfig struct {
	AllowedOrigins []string
	AllowedMethods []string
	AllowedHeaders []string
}
