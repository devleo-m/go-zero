package router

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/devleo-m/go-zero/internal/interface/http/dto"
	"github.com/devleo-m/go-zero/internal/interface/http/middleware"
)

// SetupV1Routes configura as rotas da API v1
func (r *Router) SetupV1Routes() {
	// Grupo principal da API v1
	v1 := r.engine.Group("/api/v1")

	// Middleware de versioning
	v1.Use(r.versionMiddleware("v1"))

	// Configurar sub-rotas
	r.setupUserV1Routes(v1)
	r.setupAuthV1Routes(v1)
	r.setupAdminV1Routes(v1)
	r.setupPublicV1Routes(v1)
}

// setupUserV1Routes configura rotas de usuário da v1
func (r *Router) setupUserV1Routes(v1 *gin.RouterGroup) {
	users := v1.Group("/users")

	// Rotas públicas de usuário
	users.POST("", r.userHandler.CreateUser)

	// Rotas protegidas de usuário
	protected := users.Group("")
	protected.Use(middleware.AuthMiddleware(r.jwtService, r.logger))
	{
		// Buscar usuário por ID
		protected.GET("/:id", r.userHandler.GetUser)

		// Buscar usuário por email
		protected.GET("/email/:email", r.userHandler.GetUserByEmail)

		// Listar usuários
		protected.GET("", r.userHandler.ListUsers)

		// Atualizar usuário
		protected.PUT("/:id", r.userHandler.UpdateUser)

		// Alterar senha
		protected.PUT("/:id/password", r.userHandler.ChangePassword)

		// Verificar se usuário existe
		protected.GET("/:id/exists", r.userHandler.CheckUserExists)
	}
}

// setupAuthV1Routes configura rotas de autenticação da v1
func (r *Router) setupAuthV1Routes(v1 *gin.RouterGroup) {
	auth := v1.Group("/auth")

	// Login
	auth.POST("/login", r.userHandler.AuthenticateUser)

	// Rotas protegidas de autenticação
	protected := auth.Group("")
	protected.Use(middleware.AuthMiddleware(r.jwtService, r.logger))
	{
		// Logout (se implementado)
		protected.POST("/logout", r.logout)

		// Refresh token (se implementado)
		protected.POST("/refresh", r.refreshToken)

		// Verificar token
		protected.GET("/verify", r.verifyToken)
	}
}

// setupAdminV1Routes configura rotas administrativas da v1
func (r *Router) setupAdminV1Routes(v1 *gin.RouterGroup) {
	admin := v1.Group("/admin")

	// Middleware de autenticação e autorização
	admin.Use(middleware.AuthMiddleware(r.jwtService, r.logger))
	admin.Use(middleware.RequireRole("admin"))

	// Rotas de gerenciamento de usuários
	users := admin.Group("/users")
	{
		// Ativar usuário
		users.POST("/:id/activate", r.userHandler.ActivateUser)

		// Desativar usuário
		users.POST("/:id/deactivate", r.userHandler.DeactivateUser)

		// Suspender usuário
		users.POST("/:id/suspend", r.userHandler.SuspendUser)

		// Alterar role de usuário
		users.PUT("/:id/role", r.userHandler.ChangeRole)

		// Estatísticas de usuários
		users.GET("/stats", r.userHandler.GetUserStats)
	}

	// Rotas de gerenciamento do sistema
	system := admin.Group("/system")
	{
		// Estatísticas gerais
		system.GET("/stats", r.getSystemStats)

		// Logs do sistema
		system.GET("/logs", r.getSystemLogs)

		// Configurações do sistema
		system.GET("/config", r.getSystemConfig)
	}
}

// setupPublicV1Routes configura rotas públicas da v1
func (r *Router) setupPublicV1Routes(v1 *gin.RouterGroup) {
	public := v1.Group("/public")

	// Informações públicas
	public.GET("/info", r.getPublicInfo)

	// Status da API
	public.GET("/status", r.getAPIStatus)
}

// ==========================================
// AUTH ROUTES
// ==========================================

// logout faz logout do usuário
func (r *Router) logout(c *gin.Context) {
	// Obter token do header
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse(
			"BAD_REQUEST",
			"Authorization header is required",
		))
		return
	}

	// Aqui você pode implementar a lógica de logout
	// Por exemplo, adicionar o token a uma blacklist

	response := dto.SuccessResponse{
		Success: true,
		Message: "Logged out successfully",
	}

	c.JSON(http.StatusOK, response)
}

// refreshToken renova o token de acesso
func (r *Router) refreshToken(c *gin.Context) {
	// Obter refresh token do header ou body
	refreshToken := c.GetHeader("X-Refresh-Token")
	if refreshToken == "" {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse(
			"BAD_REQUEST",
			"Refresh token is required",
		))
		return
	}

	// Aqui você pode implementar a lógica de refresh token
	// Por exemplo, validar o refresh token e gerar um novo access token

	response := dto.SuccessResponse{
		Success: true,
		Message: "Token refreshed successfully",
		Data: map[string]interface{}{
			"access_token": "new_access_token_here",
			"expires_in":   3600,
		},
	}

	c.JSON(http.StatusOK, response)
}

// verifyToken verifica se o token é válido
func (r *Router) verifyToken(c *gin.Context) {
	// O middleware de autenticação já validou o token
	// Aqui apenas retornamos informações do usuário

	userID, _ := middleware.GetUserIDFromContext(c)
	userRole, _ := middleware.GetUserRoleFromContext(c)
	userEmail, _ := middleware.GetUserEmailFromContext(c)

	response := dto.SuccessResponse{
		Success: true,
		Message: "Token is valid",
		Data: map[string]interface{}{
			"user_id":    userID,
			"user_role":  userRole,
			"user_email": userEmail,
		},
	}

	c.JSON(http.StatusOK, response)
}

// ==========================================
// ADMIN ROUTES
// ==========================================

// getSystemStats retorna estatísticas do sistema
func (r *Router) getSystemStats(c *gin.Context) {
	// Aqui você pode implementar a lógica para obter estatísticas do sistema

	response := dto.SuccessResponse{
		Success: true,
		Message: "System statistics retrieved successfully",
		Data: map[string]interface{}{
			"total_users":     1000,
			"active_users":    850,
			"pending_users":   50,
			"suspended_users": 25,
			"system_uptime":   "2 days, 5 hours",
			"memory_usage":    "512MB",
			"cpu_usage":       "15%",
		},
	}

	c.JSON(http.StatusOK, response)
}

// getSystemLogs retorna logs do sistema
func (r *Router) getSystemLogs(c *gin.Context) {
	// Aqui você pode implementar a lógica para obter logs do sistema

	response := dto.SuccessResponse{
		Success: true,
		Message: "System logs retrieved successfully",
		Data: map[string]interface{}{
			"logs": []map[string]interface{}{
				{
					"timestamp": "2024-01-15T10:30:00Z",
					"level":     "INFO",
					"message":   "User created successfully",
					"user_id":   "550e8400-e29b-41d4-a716-446655440000",
				},
				{
					"timestamp": "2024-01-15T10:31:00Z",
					"level":     "WARN",
					"message":   "Failed login attempt",
					"ip":        "192.168.1.100",
				},
			},
		},
	}

	c.JSON(http.StatusOK, response)
}

// getSystemConfig retorna configurações do sistema
func (r *Router) getSystemConfig(c *gin.Context) {
	// Aqui você pode implementar a lógica para obter configurações do sistema

	response := dto.SuccessResponse{
		Success: true,
		Message: "System configuration retrieved successfully",
		Data: map[string]interface{}{
			"app_name":    "GO ZERO",
			"version":     "1.0.0",
			"environment": "development",
			"features": map[string]bool{
				"user_registration":  true,
				"email_verification": true,
				"password_reset":     true,
				"two_factor_auth":    false,
			},
		},
	}

	c.JSON(http.StatusOK, response)
}

// ==========================================
// PUBLIC ROUTES
// ==========================================

// getPublicInfo retorna informações públicas da API
func (r *Router) getPublicInfo(c *gin.Context) {
	response := dto.SuccessResponse{
		Success: true,
		Message: "Public information retrieved successfully",
		Data: map[string]interface{}{
			"app_name":    "GO ZERO",
			"version":     "1.0.0",
			"description": "A comprehensive learning project for Go backend development",
			"features": []string{
				"User Management",
				"Authentication",
				"Role-based Access Control",
				"RESTful API",
				"Hexagonal Architecture",
			},
			"contact": map[string]string{
				"email":   "contact@gozero.dev",
				"website": "https://gozero.dev",
			},
		},
	}

	c.JSON(http.StatusOK, response)
}

// ==========================================
// MIDDLEWARE
// ==========================================

// versionMiddleware adiciona informações de versão ao contexto
func (r *Router) versionMiddleware(version string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("API-Version", version)
		c.Set("api_version", version)
		c.Next()
	}
}
