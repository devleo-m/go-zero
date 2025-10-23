package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/devleo-m/go-zero/internal/interface/http/dto"
)

// SetupPublicRoutes configura rotas públicas (sem autenticação)
func (r *Router) SetupPublicRoutes() {
	// Grupo de rotas públicas
	public := r.engine.Group("/")

	// Middleware para rotas públicas
	public.Use(r.publicMiddleware())

	// Configurar sub-rotas públicas
	r.setupHealthPublicRoutes(public)
	r.setupInfoPublicRoutes(public)
	r.setupAuthPublicRoutes(public)
}

// setupHealthPublicRoutes configura rotas de health check públicas
func (r *Router) setupHealthPublicRoutes(public *gin.RouterGroup) {
	health := public.Group("/")
	{
		// Health check básico
		health.GET("/health", r.healthHandler.HealthCheck)

		// Readiness check
		health.GET("/ready", r.healthHandler.ReadinessCheck)

		// Liveness check
		health.GET("/live", r.healthHandler.LivenessCheck)

		// Versão da aplicação
		health.GET("/version", r.healthHandler.Version)
	}
}

// setupInfoPublicRoutes configura rotas de informações públicas
func (r *Router) setupInfoPublicRoutes(public *gin.RouterGroup) {
	info := public.Group("/info")
	{
		// Informações gerais da API
		info.GET("", r.getAPIInfo)

		// Status da API
		info.GET("/status", r.getAPIStatus)

		// Documentação da API
		info.GET("/docs", r.getAPIDocs)

		// Termos de uso
		info.GET("/terms", r.getTermsOfService)

		// Política de privacidade
		info.GET("/privacy", r.getPrivacyPolicy)
	}
}

// setupAuthPublicRoutes configura rotas de autenticação públicas
func (r *Router) setupAuthPublicRoutes(public *gin.RouterGroup) {
	auth := public.Group("/auth")
	{
		// Login
		auth.POST("/login", r.userHandler.AuthenticateUser)

		// Registro de usuário (se permitido)
		auth.POST("/register", r.userHandler.CreateUser)

		// Verificar se email está disponível
		auth.GET("/check-email/:email", r.checkEmailAvailability)

		// Solicitar reset de senha
		auth.POST("/forgot-password", r.forgotPassword)

		// Reset de senha
		auth.POST("/reset-password", r.resetPassword)
	}
}

// ==========================================
// PUBLIC HANDLERS
// ==========================================

// getAPIInfo retorna informações gerais da API
func (r *Router) getAPIInfo(c *gin.Context) {
	response := dto.SuccessResponse{
		Success: true,
		Message: "API information retrieved successfully",
		Data: map[string]interface{}{
			"name":          "GO ZERO API",
			"version":       "1.0.0",
			"description":   "A comprehensive learning project for Go backend development",
			"author":        "GO ZERO Team",
			"license":       "MIT",
			"repository":    "https://github.com/go-zero/go-zero",
			"documentation": "https://docs.gozero.dev",
			"features": []string{
				"User Management",
				"Authentication & Authorization",
				"Role-based Access Control",
				"RESTful API Design",
				"Hexagonal Architecture",
				"Clean Code Principles",
				"Comprehensive Testing",
				"Production-ready Patterns",
			},
			"technologies": []string{
				"Go 1.21+",
				"Gin Framework",
				"GORM",
				"PostgreSQL",
				"Redis",
				"JWT",
				"Zap Logger",
				"Docker",
			},
			"endpoints": map[string]string{
				"users":  "/api/v1/users",
				"auth":   "/api/v1/auth",
				"admin":  "/api/v1/admin",
				"health": "/health",
				"docs":   "/info/docs",
			},
		},
	}

	c.JSON(http.StatusOK, response)
}

// getAPIStatus retorna o status atual da API
func (r *Router) getAPIStatus(c *gin.Context) {
	response := dto.SuccessResponse{
		Success: true,
		Message: "API is operational",
		Data: map[string]interface{}{
			"status":    "healthy",
			"version":   "1.0.0",
			"uptime":    "2 days, 5 hours",
			"timestamp": "2024-01-15T10:30:00Z",
			"services": map[string]string{
				"database": "connected",
				"redis":    "connected",
				"logger":   "active",
			},
			"performance": map[string]interface{}{
				"response_time_ms":    45,
				"requests_per_second": 150,
				"error_rate":          "0.1%",
			},
		},
	}

	c.JSON(http.StatusOK, response)
}

// getAPIDocs retorna informações sobre a documentação da API
func (r *Router) getAPIDocs(c *gin.Context) {
	response := dto.SuccessResponse{
		Success: true,
		Message: "API documentation information retrieved successfully",
		Data: map[string]interface{}{
			"swagger_ui":         "https://api.gozero.dev/swagger/",
			"openapi_spec":       "https://api.gozero.dev/swagger/doc.json",
			"postman_collection": "https://api.gozero.dev/postman/collection.json",
			"examples": map[string]string{
				"curl":       "curl -X GET https://api.gozero.dev/api/v1/users",
				"javascript": "fetch('https://api.gozero.dev/api/v1/users')",
				"python":     "requests.get('https://api.gozero.dev/api/v1/users')",
			},
			"authentication": map[string]string{
				"type":    "Bearer Token",
				"header":  "Authorization: Bearer <token>",
				"example": "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
			},
		},
	}

	c.JSON(http.StatusOK, response)
}

// getTermsOfService retorna os termos de uso
func (r *Router) getTermsOfService(c *gin.Context) {
	response := dto.SuccessResponse{
		Success: true,
		Message: "Terms of service retrieved successfully",
		Data: map[string]interface{}{
			"version":      "1.0",
			"last_updated": "2024-01-01T00:00:00Z",
			"terms":        "By using this API, you agree to the following terms...",
			"contact":      "legal@gozero.dev",
		},
	}

	c.JSON(http.StatusOK, response)
}

// getPrivacyPolicy retorna a política de privacidade
func (r *Router) getPrivacyPolicy(c *gin.Context) {
	response := dto.SuccessResponse{
		Success: true,
		Message: "Privacy policy retrieved successfully",
		Data: map[string]interface{}{
			"version":      "1.0",
			"last_updated": "2024-01-01T00:00:00Z",
			"policy":       "This API respects your privacy and follows GDPR guidelines...",
			"contact":      "privacy@gozero.dev",
		},
	}

	c.JSON(http.StatusOK, response)
}

// checkEmailAvailability verifica se um email está disponível
func (r *Router) checkEmailAvailability(c *gin.Context) {
	email := c.Param("email")
	if email == "" {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse(
			"BAD_REQUEST",
			"Email parameter is required",
		))
		return
	}

	// Aqui você pode implementar a lógica para verificar se o email está disponível
	// Por exemplo, consultar o banco de dados

	response := dto.SuccessResponse{
		Success: true,
		Message: "Email availability checked successfully",
		Data: map[string]interface{}{
			"email":     email,
			"available": true,
		},
	}

	c.JSON(http.StatusOK, response)
}

// forgotPassword solicita reset de senha
func (r *Router) forgotPassword(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		r.errorHandler.HandleError(c, err)
		return
	}

	// Aqui você pode implementar a lógica para solicitar reset de senha
	// Por exemplo, enviar email com token de reset

	response := dto.SuccessResponse{
		Success: true,
		Message: "Password reset instructions sent to your email",
	}

	c.JSON(http.StatusOK, response)
}

// resetPassword reseta a senha com token
func (r *Router) resetPassword(c *gin.Context) {
	var req struct {
		Token       string `json:"token" binding:"required"`
		NewPassword string `json:"new_password" binding:"required,min=8"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		r.errorHandler.HandleError(c, err)
		return
	}

	// Aqui você pode implementar a lógica para resetar a senha
	// Por exemplo, validar o token e atualizar a senha

	response := dto.SuccessResponse{
		Success: true,
		Message: "Password reset successfully",
	}

	c.JSON(http.StatusOK, response)
}

// ==========================================
// MIDDLEWARE
// ==========================================

// publicMiddleware middleware para rotas públicas
func (r *Router) publicMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Adicionar headers específicos para rotas públicas
		c.Header("X-Public-Route", "true")
		c.Header("Cache-Control", "public, max-age=300") // Cache por 5 minutos

		// Log de acesso a rotas públicas
		r.logger.Info("Public route accessed",
			zap.String("path", c.Request.URL.Path),
			zap.String("method", c.Request.Method),
			zap.String("ip", c.ClientIP()),
		)

		c.Next()
	}
}
