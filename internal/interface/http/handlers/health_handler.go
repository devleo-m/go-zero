package handlers

import (
	"context"
	"net/http"
	"runtime"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/devleo-m/go-zero/internal/interface/http/dto"
)

// ==========================================
// INTERFACES DE DEPENDÊNCIAS
// ==========================================

// DatabaseChecker interface para verificar status do banco
type DatabaseChecker interface {
	Ping(ctx context.Context) error
}

// CacheChecker interface para verificar status do cache
type CacheChecker interface {
	Ping(ctx context.Context) error
}

// ==========================================
// HEALTH HANDLER
// ==========================================

// HealthHandler handler para health checks
type HealthHandler struct {
	logger       *zap.Logger
	db           *gorm.DB
	cache        CacheChecker
	startTime    time.Time
	version      string
	environment  string
	buildTime    string
	gitCommit    string
	mu           sync.RWMutex
	healthStatus *HealthStatus
}

// HealthStatus representa o status de saúde da aplicação
type HealthStatus struct {
	Status      string                     `json:"status"`
	Timestamp   time.Time                  `json:"timestamp"`
	Components  map[string]ComponentHealth `json:"components"`
	Version     string                     `json:"version"`
	Uptime      string                     `json:"uptime"`
	Environment string                     `json:"environment"`
}

// ComponentHealth representa o status de um componente
type ComponentHealth struct {
	Status       string        `json:"status"`
	Message      string        `json:"message,omitempty"`
	ResponseTime time.Duration `json:"response_time_ms"`
	LastChecked  time.Time     `json:"last_checked"`
}

// HealthHandlerConfig configurações do health handler
type HealthHandlerConfig struct {
	Logger      *zap.Logger
	DB          *gorm.DB
	Cache       CacheChecker
	Version     string
	Environment string
	BuildTime   string
	GitCommit   string
}

// NewHealthHandler cria uma nova instância do HealthHandler
func NewHealthHandler(config HealthHandlerConfig) *HealthHandler {
	handler := &HealthHandler{
		logger:      config.Logger,
		db:          config.DB,
		cache:       config.Cache,
		startTime:   time.Now(),
		version:     config.Version,
		environment: config.Environment,
		buildTime:   config.BuildTime,
		gitCommit:   config.GitCommit,
		healthStatus: &HealthStatus{
			Status:     "initializing",
			Timestamp:  time.Now(),
			Components: make(map[string]ComponentHealth),
		},
	}

	// Iniciar verificação periódica de saúde
	go handler.periodicHealthCheck()

	return handler
}

// ==========================================
// HEALTH CHECK ENDPOINTS
// ==========================================

// HealthCheck verifica a saúde da aplicação
// @Summary Health check
// @Description Check if the application is healthy
// @Tags health
// @Produce json
// @Success 200 {object} dto.HealthCheckResponse
// @Failure 503 {object} dto.ErrorResponse
// @Router /health [get]
func (h *HealthHandler) HealthCheck(c *gin.Context) {
	// Obter status atual
	h.mu.RLock()
	currentStatus := h.healthStatus
	h.mu.RUnlock()

	// Verificar se está saudável
	isHealthy := currentStatus.Status == "healthy"

	response := dto.HealthCheckResponse{
		Success:   isHealthy,
		Message:   h.getHealthMessage(currentStatus.Status),
		Timestamp: time.Now(),
		Version:   h.version,
		Uptime:    time.Since(h.startTime).String(),
	}

	statusCode := http.StatusOK
	if !isHealthy {
		statusCode = http.StatusServiceUnavailable
		h.logger.Warn("Health check failed",
			zap.String("status", currentStatus.Status),
			zap.Any("components", currentStatus.Components),
		)
	}

	c.JSON(statusCode, response)
}

// ReadinessCheck verifica se a aplicação está pronta para receber tráfego
// @Summary Readiness check
// @Description Check if the application is ready to receive traffic
// @Tags health
// @Produce json
// @Success 200 {object} dto.HealthCheckResponse
// @Failure 503 {object} dto.ErrorResponse
// @Router /ready [get]
func (h *HealthHandler) ReadinessCheck(c *gin.Context) {
	ctx := c.Request.Context()

	// Verificar todas as dependências críticas
	components := h.checkAllComponents(ctx)

	// Verificar se todos os componentes críticos estão OK
	isReady := h.areComponentsHealthy(components, []string{"database"})

	response := dto.HealthCheckResponse{
		Success:   isReady,
		Message:   h.getReadinessMessage(isReady),
		Timestamp: time.Now(),
		Version:   h.version,
		Uptime:    time.Since(h.startTime).String(),
	}

	statusCode := http.StatusOK
	if !isReady {
		statusCode = http.StatusServiceUnavailable
		h.logger.Warn("Readiness check failed",
			zap.Any("components", components),
		)
	}

	c.JSON(statusCode, response)
}

// LivenessCheck verifica se a aplicação está viva
// @Summary Liveness check
// @Description Check if the application is alive
// @Tags health
// @Produce json
// @Success 200 {object} dto.HealthCheckResponse
// @Router /live [get]
func (h *HealthHandler) LivenessCheck(c *gin.Context) {
	// Liveness check é simples: se conseguir responder, está vivo
	uptime := time.Since(h.startTime)
	response := dto.HealthCheckResponse{
		Success:   true,
		Message:   "Service is alive",
		Timestamp: time.Now(),
		Version:   h.version,
		Uptime:    uptime.String(),
	}

	c.JSON(http.StatusOK, response)
}

// ==========================================
// VERSION ENDPOINT
// ==========================================

// Version retorna a versão da aplicação
// @Summary Get application version
// @Description Get the current version of the application
// @Tags health
// @Produce json
// @Success 200 {object} dto.SuccessResponse
// @Router /version [get]
func (h *HealthHandler) Version(c *gin.Context) {
	response := dto.SuccessResponse{
		Success: true,
		Message: "Version retrieved successfully",
		Data: map[string]interface{}{
			"version":     h.version,
			"build_time":  h.buildTime,
			"git_commit":  h.gitCommit,
			"go_version":  runtime.Version(),
			"environment": h.environment,
			"start_time":  h.startTime.Format(time.RFC3339),
			"uptime":      time.Since(h.startTime).String(),
		},
	}

	c.JSON(http.StatusOK, response)
}

// ==========================================
// METRICS ENDPOINT
// ==========================================

// Metrics retorna métricas básicas da aplicação
// @Summary Get application metrics
// @Description Get basic application metrics
// @Tags health
// @Produce json
// @Success 200 {object} dto.SuccessResponse
// @Security BearerAuth
// @Router /metrics [get]
func (h *HealthHandler) Metrics(c *gin.Context) {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	uptime := time.Since(h.startTime)

	response := dto.SuccessResponse{
		Success: true,
		Message: "Metrics retrieved successfully",
		Data: map[string]interface{}{
			// Uptime
			"uptime":         uptime.String(),
			"uptime_seconds": int64(uptime.Seconds()),
			"start_time":     h.startTime.Format(time.RFC3339),

			// Version info
			"version":     h.version,
			"environment": h.environment,
			"go_version":  runtime.Version(),

			// Memory
			"memory": map[string]interface{}{
				"alloc_mb":       memStats.Alloc / 1024 / 1024,
				"total_alloc_mb": memStats.TotalAlloc / 1024 / 1024,
				"sys_mb":         memStats.Sys / 1024 / 1024,
				"num_gc":         memStats.NumGC,
			},

			// Runtime
			"runtime": map[string]interface{}{
				"num_goroutine": runtime.NumGoroutine(),
				"num_cpu":       runtime.NumCPU(),
				"gomaxprocs":    runtime.GOMAXPROCS(0),
			},

			// Status
			"status": h.healthStatus.Status,
		},
	}

	c.JSON(http.StatusOK, response)
}

// DetailedHealthCheck retorna status detalhado de saúde
// @Summary Detailed health check
// @Description Get detailed health status of all components
// @Tags health
// @Produce json
// @Success 200 {object} dto.SuccessResponse
// @Failure 503 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /health/detailed [get]
func (h *HealthHandler) DetailedHealthCheck(c *gin.Context) {
	ctx := c.Request.Context()

	// Verificar todos os componentes
	components := h.checkAllComponents(ctx)

	// Determinar status geral
	overallStatus := "healthy"
	if !h.areComponentsHealthy(components, []string{"database"}) {
		overallStatus = "unhealthy"
	} else if !h.areComponentsHealthy(components, []string{"cache"}) {
		overallStatus = "degraded"
	}

	healthStatus := HealthStatus{
		Status:      overallStatus,
		Timestamp:   time.Now(),
		Components:  components,
		Version:     h.version,
		Uptime:      time.Since(h.startTime).String(),
		Environment: h.environment,
	}

	response := dto.SuccessResponse{
		Success: overallStatus == "healthy",
		Message: h.getHealthMessage(overallStatus),
		Data:    healthStatus,
	}

	statusCode := http.StatusOK
	if overallStatus == "unhealthy" {
		statusCode = http.StatusServiceUnavailable
	}

	c.JSON(statusCode, response)
}

// ==========================================
// COMPONENT HEALTH CHECKS
// ==========================================

// checkAllComponents verifica todos os componentes
func (h *HealthHandler) checkAllComponents(ctx context.Context) map[string]ComponentHealth {
	components := make(map[string]ComponentHealth)

	// Verificar banco de dados
	components["database"] = h.checkDatabase(ctx)

	// Verificar cache (se configurado)
	if h.cache != nil {
		components["cache"] = h.checkCache(ctx)
	}

	return components
}

// checkDatabase verifica a saúde do banco de dados
func (h *HealthHandler) checkDatabase(ctx context.Context) ComponentHealth {
	start := time.Now()

	// Verificar conexão com o banco
	sqlDB, err := h.db.DB()
	if err != nil {
		return ComponentHealth{
			Status:       "unhealthy",
			Message:      "Failed to get database connection: " + err.Error(),
			ResponseTime: time.Since(start),
			LastChecked:  time.Now(),
		}
	}

	// Ping no banco
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		return ComponentHealth{
			Status:       "unhealthy",
			Message:      "Database ping failed: " + err.Error(),
			ResponseTime: time.Since(start),
			LastChecked:  time.Now(),
		}
	}

	// Verificar estatísticas do pool
	stats := sqlDB.Stats()
	if stats.OpenConnections == 0 {
		return ComponentHealth{
			Status:       "unhealthy",
			Message:      "No open database connections",
			ResponseTime: time.Since(start),
			LastChecked:  time.Now(),
		}
	}

	return ComponentHealth{
		Status:       "healthy",
		Message:      "Database is operational",
		ResponseTime: time.Since(start),
		LastChecked:  time.Now(),
	}
}

// checkCache verifica a saúde do cache
func (h *HealthHandler) checkCache(ctx context.Context) ComponentHealth {
	start := time.Now()
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	if err := h.cache.Ping(ctx); err != nil {
		return ComponentHealth{
			Status:       "degraded",
			Message:      "Cache ping failed: " + err.Error(),
			ResponseTime: time.Since(start),
			LastChecked:  time.Now(),
		}
	}

	return ComponentHealth{
		Status:       "healthy",
		Message:      "Cache is operational",
		ResponseTime: time.Since(start),
		LastChecked:  time.Now(),
	}
}

// ==========================================
// HELPER METHODS
// ==========================================

// areComponentsHealthy verifica se os componentes especificados estão saudáveis
func (h *HealthHandler) areComponentsHealthy(components map[string]ComponentHealth, criticalComponents []string) bool {
	for _, name := range criticalComponents {
		if component, exists := components[name]; exists {
			if component.Status == "unhealthy" {
				return false
			}
		}
	}
	return true
}

// getHealthMessage retorna mensagem apropriada para o status
func (h *HealthHandler) getHealthMessage(status string) string {
	switch status {
	case "healthy":
		return "Service is healthy"
	case "degraded":
		return "Service is operational but degraded"
	case "unhealthy":
		return "Service is unhealthy"
	default:
		return "Service status unknown"
	}
}

// getReadinessMessage retorna mensagem de readiness
func (h *HealthHandler) getReadinessMessage(isReady bool) string {
	if isReady {
		return "Service is ready to receive traffic"
	}
	return "Service is not ready to receive traffic"
}

// periodicHealthCheck executa verificação periódica de saúde
func (h *HealthHandler) periodicHealthCheck() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		components := h.checkAllComponents(ctx)

		// Determinar status geral
		status := "healthy"
		if !h.areComponentsHealthy(components, []string{"database"}) {
			status = "unhealthy"
		} else if !h.areComponentsHealthy(components, []string{"cache"}) {
			status = "degraded"
		}

		// Atualizar status
		h.mu.Lock()
		h.healthStatus = &HealthStatus{
			Status:      status,
			Timestamp:   time.Now(),
			Components:  components,
			Version:     h.version,
			Uptime:      time.Since(h.startTime).String(),
			Environment: h.environment,
		}
		h.mu.Unlock()

		cancel()

		// Log se não estiver saudável
		if status != "healthy" {
			h.logger.Warn("Periodic health check detected issues",
				zap.String("status", status),
				zap.Any("components", components),
			)
		}
	}
}
