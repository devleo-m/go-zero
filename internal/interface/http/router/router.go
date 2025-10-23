package router

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/devleo-m/go-zero/internal/interface/http/dto"
	"github.com/devleo-m/go-zero/internal/interface/http/handlers"
	"github.com/devleo-m/go-zero/internal/interface/http/middleware"
	"github.com/devleo-m/go-zero/internal/interface/http/validation"
	"github.com/devleo-m/go-zero/internal/usecase/user"
)

// Router configuração do router HTTP
type Router struct {
	engine        *gin.Engine
	userHandler   *handlers.UserHandler
	healthHandler *handlers.HealthHandler
	errorHandler  *handlers.ErrorHandler
	validator     *validation.CustomValidator
	logger        *zap.Logger
	jwtService    JWTService
}

// JWTService interface para serviços de JWT
type JWTService interface {
	ValidateToken(token string) (map[string]interface{}, error)
	ExtractUserID(claims map[string]interface{}) (string, error)
	ExtractUserRole(claims map[string]interface{}) (string, error)
	ExtractUserEmail(claims map[string]interface{}) (string, error)
}

// Config configuração do router
type Config struct {
	UserUseCaseAggregate user.UserUseCaseAggregate
	JWTService           JWTService
	Logger               *zap.Logger
	Environment          string
	EnableSwagger        bool
	EnableMetrics        bool
}

// NewRouter cria uma nova instância do Router
func NewRouter(config Config) *Router {
	// Configurar modo do Gin baseado no ambiente
	switch config.Environment {
	case "production":
		gin.SetMode(gin.ReleaseMode)
	case "test":
		gin.SetMode(gin.TestMode)
	default:
		gin.SetMode(gin.DebugMode)
	}

	// Criar engine do Gin
	engine := gin.New()

	// Criar instâncias dos handlers
	validator := validation.NewCustomValidator()
	errorHandler := handlers.NewErrorHandler(config.Logger)
	userHandler := handlers.NewUserHandler(
		config.UserUseCaseAggregate,
		validator,
		errorHandler,
		config.Logger,
	)
	healthHandler := handlers.NewHealthHandler(config.Logger)

	router := &Router{
		engine:        engine,
		userHandler:   userHandler,
		healthHandler: healthHandler,
		errorHandler:  errorHandler,
		validator:     validator,
		logger:        config.Logger,
		jwtService:    config.JWTService,
	}

	// Configurar middlewares globais
	router.setupGlobalMiddlewares()

	// Configurar rotas
	router.setupRoutes()

	// Configurar Swagger se habilitado
	if config.EnableSwagger {
		router.setupSwagger()
	}

	// Configurar métricas se habilitado
	if config.EnableMetrics {
		router.setupMetrics()
	}

	return router
}

// setupGlobalMiddlewares configura middlewares globais
func (r *Router) setupGlobalMiddlewares() {
	// Middleware de recovery (deve ser o primeiro)
	r.engine.Use(middleware.RecoveryMiddleware(r.logger))

	// Middleware de CORS
	r.engine.Use(middleware.CORSForDevelopment())

	// Middleware de request ID
	r.engine.Use(middleware.RequestIDMiddleware())

	// Middleware de logging
	r.engine.Use(middleware.LoggerMiddleware(r.logger))

	// Middleware de headers de segurança
	r.engine.Use(middleware.SecurityHeadersMiddleware())

	// Middleware de rate limiting
	r.engine.Use(middleware.RateLimitMiddleware())

	// Middleware de timeout
	r.engine.Use(middleware.TimeoutMiddleware(30 * time.Second))

	// Middleware de tamanho máximo do body
	r.engine.Use(middleware.MaxBodySizeMiddleware(10 * 1024 * 1024)) // 10MB

	// Middleware de tratamento de erros
	r.engine.Use(r.errorHandler.ErrorHandlerMiddleware())
}

// setupRoutes configura todas as rotas
func (r *Router) setupRoutes() {
	// Rotas de health check (sem autenticação)
	r.setupHealthRoutes()

	// Rotas públicas (sem autenticação)
	r.setupPublicRoutes()

	// Rotas protegidas (com autenticação)
	r.setupProtectedRoutes()

	// Rotas administrativas (com autenticação e permissões)
	r.setupAdminRoutes()

	// Rota de fallback para 404
	r.engine.NoRoute(r.handleNotFound)
}

// setupHealthRoutes configura rotas de health check
func (r *Router) setupHealthRoutes() {
	health := r.engine.Group("/")
	{
		health.GET("/health", r.healthHandler.HealthCheck)
		health.GET("/ready", r.healthHandler.ReadinessCheck)
		health.GET("/live", r.healthHandler.LivenessCheck)
		health.GET("/version", r.healthHandler.Version)
	}
}

// setupPublicRoutes configura rotas públicas
func (r *Router) setupPublicRoutes() {
	public := r.engine.Group("/api/v1")
	{
		// Autenticação
		auth := public.Group("/auth")
		{
			auth.POST("/login", r.userHandler.AuthenticateUser)
		}
	}
}

// setupProtectedRoutes configura rotas protegidas
func (r *Router) setupProtectedRoutes() {
	// Grupo de rotas protegidas
	protected := r.engine.Group("/api/v1")
	protected.Use(middleware.AuthMiddleware(r.jwtService, r.logger))
	{
		// Rotas de usuário
		users := protected.Group("/users")
		{
			// Criar usuário (requer autenticação)
			users.POST("", r.userHandler.CreateUser)

			// Buscar usuário por ID
			users.GET("/:id", r.userHandler.GetUser)

			// Buscar usuário por email
			users.GET("/email/:email", r.userHandler.GetUserByEmail)

			// Listar usuários
			users.GET("", r.userHandler.ListUsers)

			// Atualizar usuário
			users.PUT("/:id", r.userHandler.UpdateUser)

			// Alterar senha
			users.PUT("/:id/password", r.userHandler.ChangePassword)

			// Verificar se usuário existe
			users.GET("/:id/exists", r.userHandler.CheckUserExists)
		}

		// Rotas de perfil do usuário logado
		profile := protected.Group("/profile")
		{
			// Buscar perfil do usuário logado
			profile.GET("", r.getCurrentUserProfile)

			// Atualizar perfil do usuário logado
			profile.PUT("", r.updateCurrentUserProfile)

			// Alterar senha do usuário logado
			profile.PUT("/password", r.changeCurrentUserPassword)
		}
	}
}

// setupAdminRoutes configura rotas administrativas
func (r *Router) setupAdminRoutes() {
	// Grupo de rotas administrativas
	admin := r.engine.Group("/api/v1/admin")
	admin.Use(middleware.AuthMiddleware(r.jwtService, r.logger))
	admin.Use(middleware.RequireRole("admin"))
	{
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
	}
}

// setupSwagger configura Swagger/OpenAPI
func (r *Router) setupSwagger() {
	// Aqui você pode adicionar configuração do Swagger
	// Por exemplo, usando gin-swagger
	r.logger.Info("Swagger documentation enabled")
}

// setupMetrics configura métricas
func (r *Router) setupMetrics() {
	// Grupo de rotas de métricas
	metrics := r.engine.Group("/metrics")
	metrics.Use(middleware.AuthMiddleware(r.jwtService, r.logger))
	{
		metrics.GET("", r.healthHandler.Metrics)
	}
}

// handleNotFound trata requisições para rotas não encontradas
func (r *Router) handleNotFound(c *gin.Context) {
	response := dto.NewErrorResponse(
		"NOT_FOUND",
		"The requested resource was not found",
	)
	c.JSON(http.StatusNotFound, response)
}

// ==========================================
// PROFILE ROUTES (CURRENT USER)
// ==========================================

// getCurrentUserProfile busca o perfil do usuário logado
func (r *Router) getCurrentUserProfile(c *gin.Context) {
	// Obter ID do usuário logado
	userID, exists := middleware.GetUserIDFromContext(c)
	if !exists {
		r.errorHandler.HandleCustomError(c, "UNAUTHORIZED", "User not authenticated", http.StatusUnauthorized)
		return
	}

	// Buscar usuário
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		r.errorHandler.HandleValidationError(c, "user_id", "Invalid user ID format")
		return
	}

	input := user.GetUserInput{
		ID: userUUID,
	}

	output, err := r.userHandler.GetUseCaseAggregate().GetUser(c.Request.Context(), input)
	if err != nil {
		r.errorHandler.HandleError(c, err)
		return
	}

	response := dto.GetUserResponse{
		Success: true,
		Message: "Profile retrieved successfully",
		Data:    dto.ToUserResponseFromOutput(output.User),
	}

	c.JSON(http.StatusOK, response)
}

// updateCurrentUserProfile atualiza o perfil do usuário logado
func (r *Router) updateCurrentUserProfile(c *gin.Context) {
	// Obter ID do usuário logado
	userID, exists := middleware.GetUserIDFromContext(c)
	if !exists {
		r.errorHandler.HandleCustomError(c, "UNAUTHORIZED", "User not authenticated", http.StatusUnauthorized)
		return
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		r.errorHandler.HandleValidationError(c, "user_id", "Invalid user ID format")
		return
	}

	var req dto.UpdateUserRequest

	// Bind e validar request
	if err := c.ShouldBindJSON(&req); err != nil {
		r.errorHandler.HandleError(c, err)
		return
	}

	// Validar com validador customizado
	if err := r.validator.ValidateStruct(req); err != nil {
		r.errorHandler.HandleError(c, err)
		return
	}

	// Converter para input do use case
	input := req.ToUpdateUserInput(userUUID)

	// Executar use case
	output, err := r.userHandler.GetUseCaseAggregate().UpdateUser(c.Request.Context(), input)
	if err != nil {
		r.errorHandler.HandleError(c, err)
		return
	}

	// Converter para response
	response := dto.UpdateUserResponse{
		Success: true,
		Message: "Profile updated successfully",
		Data:    dto.ToUserResponseFromOutput(output.User),
	}

	c.JSON(http.StatusOK, response)
}

// changeCurrentUserPassword altera a senha do usuário logado
func (r *Router) changeCurrentUserPassword(c *gin.Context) {
	// Obter ID do usuário logado
	userID, exists := middleware.GetUserIDFromContext(c)
	if !exists {
		r.errorHandler.HandleCustomError(c, "UNAUTHORIZED", "User not authenticated", http.StatusUnauthorized)
		return
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		r.errorHandler.HandleValidationError(c, "user_id", "Invalid user ID format")
		return
	}

	var req dto.ChangePasswordRequest

	// Bind e validar request
	if err := c.ShouldBindJSON(&req); err != nil {
		r.errorHandler.HandleError(c, err)
		return
	}

	// Validar com validador customizado
	if err := r.validator.ValidateStruct(req); err != nil {
		r.errorHandler.HandleError(c, err)
		return
	}

	// Verificar se as senhas coincidem
	if req.NewPassword != req.ConfirmNewPassword {
		r.errorHandler.HandleValidationError(c, "confirm_new_password", "Passwords do not match")
		return
	}

	// Converter para input do use case
	input := req.ToChangePasswordInput(userUUID)

	// Executar use case
	_, err = r.userHandler.GetUseCaseAggregate().ChangePassword(c.Request.Context(), input)
	if err != nil {
		r.errorHandler.HandleError(c, err)
		return
	}

	// Converter para response
	response := dto.ChangePasswordResponse{
		Success: true,
		Message: "Password changed successfully",
	}

	c.JSON(http.StatusOK, response)
}

// ==========================================
// ROUTER METHODS
// ==========================================

// Setup retorna o engine do Gin configurado
func (r *Router) Setup() *gin.Engine {
	return r.engine
}

// Run inicia o servidor HTTP
func (r *Router) Run(addr string) error {
	r.logger.Info("Starting HTTP server", zap.String("address", addr))
	return r.engine.Run(addr)
}

// RunTLS inicia o servidor HTTPS
func (r *Router) RunTLS(addr, certFile, keyFile string) error {
	r.logger.Info("Starting HTTPS server", zap.String("address", addr))
	return r.engine.RunTLS(addr, certFile, keyFile)
}

// ==========================================
// ROUTE GROUPS
// ==========================================

// GetUserRoutes retorna o grupo de rotas de usuário
func (r *Router) GetUserRoutes() *gin.RouterGroup {
	return r.engine.Group("/api/v1/users")
}

// GetAuthRoutes retorna o grupo de rotas de autenticação
func (r *Router) GetAuthRoutes() *gin.RouterGroup {
	return r.engine.Group("/api/v1/auth")
}

// GetAdminRoutes retorna o grupo de rotas administrativas
func (r *Router) GetAdminRoutes() *gin.RouterGroup {
	return r.engine.Group("/api/v1/admin")
}

// GetHealthRoutes retorna o grupo de rotas de health check
func (r *Router) GetHealthRoutes() *gin.RouterGroup {
	return r.engine.Group("/")
}
