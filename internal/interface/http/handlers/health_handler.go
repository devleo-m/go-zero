package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/devleo-m/go-zero/internal/interface/http/dto"
)

// HealthHandler handler para health checks
type HealthHandler struct {
	logger    *zap.Logger
	startTime time.Time
}

// NewHealthHandler cria uma nova instância do HealthHandler
func NewHealthHandler(logger *zap.Logger) *HealthHandler {
	return &HealthHandler{
		logger:    logger,
		startTime: time.Now(),
	}
}

// HealthCheck verifica a saúde da aplicação
// @Summary Health check
// @Description Check if the application is healthy
// @Tags health
// @Produce json
// @Success 200 {object} dto.HealthCheckResponse
// @Router /health [get]
func (h *HealthHandler) HealthCheck(c *gin.Context) {
	uptime := time.Since(h.startTime)

	response := dto.HealthCheckResponse{
		Success:   true,
		Message:   "Service is healthy",
		Timestamp: time.Now(),
		Version:   "1.0.0",
		Uptime:    uptime.String(),
	}

	c.JSON(http.StatusOK, response)
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
	// Aqui você pode adicionar verificações de dependências
	// como banco de dados, Redis, etc.

	uptime := time.Since(h.startTime)

	response := dto.HealthCheckResponse{
		Success:   true,
		Message:   "Service is ready",
		Timestamp: time.Now(),
		Version:   "1.0.0",
		Uptime:    uptime.String(),
	}

	c.JSON(http.StatusOK, response)
}

// LivenessCheck verifica se a aplicação está viva
// @Summary Liveness check
// @Description Check if the application is alive
// @Tags health
// @Produce json
// @Success 200 {object} dto.HealthCheckResponse
// @Router /live [get]
func (h *HealthHandler) LivenessCheck(c *gin.Context) {
	uptime := time.Since(h.startTime)

	response := dto.HealthCheckResponse{
		Success:   true,
		Message:   "Service is alive",
		Timestamp: time.Now(),
		Version:   "1.0.0",
		Uptime:    uptime.String(),
	}

	c.JSON(http.StatusOK, response)
}

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
			"version":     "1.0.0",
			"build_time":  h.startTime.Format(time.RFC3339),
			"go_version":  "1.21+",
			"environment": "development",
		},
	}

	c.JSON(http.StatusOK, response)
}

// Metrics retorna métricas básicas da aplicação
// @Summary Get application metrics
// @Description Get basic application metrics
// @Tags health
// @Produce json
// @Success 200 {object} dto.SuccessResponse
// @Security BearerAuth
// @Router /metrics [get]
func (h *HealthHandler) Metrics(c *gin.Context) {
	uptime := time.Since(h.startTime)

	response := dto.SuccessResponse{
		Success: true,
		Message: "Metrics retrieved successfully",
		Data: map[string]interface{}{
			"uptime":         uptime.String(),
			"uptime_seconds": int64(uptime.Seconds()),
			"start_time":     h.startTime.Format(time.RFC3339),
			"version":        "1.0.0",
			"status":         "healthy",
		},
	}

	c.JSON(http.StatusOK, response)
}
