package middleware

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// MetricsMiddleware middleware para coletar métricas de performance
func MetricsMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		// Processar requisição
		c.Next()

		// Calcular métricas
		latency := time.Since(start)
		status := c.Writer.Status()
		size := c.Writer.Size()

		// Log das métricas
		logger.Info("Request metrics",
			zap.String("method", method),
			zap.String("path", path),
			zap.Int("status", status),
			zap.Int64("latency_ms", latency.Milliseconds()),
			zap.Int("size", size),
			zap.String("ip", c.ClientIP()),
			zap.String("user_agent", c.Request.UserAgent()),
		)

		// Adicionar headers de métricas
		c.Header("X-Response-Time", strconv.FormatInt(latency.Milliseconds(), 10)+"ms")
		c.Header("X-Content-Length", strconv.Itoa(size))
	}
}

// SlowRequestMiddleware middleware para detectar requisições lentas
func SlowRequestMiddleware(threshold time.Duration, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		c.Next()

		latency := time.Since(start)

		// Verificar se a requisição foi lenta
		if latency > threshold {
			logger.Warn("Slow request detected",
				zap.String("method", method),
				zap.String("path", path),
				zap.Duration("latency", latency),
				zap.Duration("threshold", threshold),
				zap.String("ip", c.ClientIP()),
				zap.String("user_agent", c.Request.UserAgent()),
			)
		}
	}
}

// RequestSizeMiddleware middleware para limitar tamanho das requisições
func RequestSizeMiddleware(maxSize int64, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Verificar Content-Length
		if c.Request.ContentLength > maxSize {
			logger.Warn("Request too large",
				zap.String("path", c.Request.URL.Path),
				zap.String("method", c.Request.Method),
				zap.Int64("content_length", c.Request.ContentLength),
				zap.Int64("max_size", maxSize),
				zap.String("ip", c.ClientIP()),
			)

			c.JSON(413, gin.H{
				"success": false,
				"error":   "REQUEST_TOO_LARGE",
				"message": "Request body too large",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// ==========================================
// CUSTOM METRICS COLLECTOR
// ==========================================

// MetricsCollector interface para coletar métricas customizadas
type MetricsCollector interface {
	IncrementCounter(name string, labels map[string]string)
	RecordHistogram(name string, value float64, labels map[string]string)
	RecordGauge(name string, value float64, labels map[string]string)
}

// PrometheusMetricsMiddleware middleware para métricas do Prometheus
func PrometheusMetricsMiddleware(collector MetricsCollector) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		c.Next()

		// Calcular métricas
		latency := time.Since(start).Seconds()
		status := c.Writer.Status()
		size := c.Writer.Size()

		// Labels para as métricas
		labels := map[string]string{
			"method": method,
			"path":   path,
			"status": strconv.Itoa(status),
		}

		// Incrementar contador de requisições
		collector.IncrementCounter("http_requests_total", labels)

		// Registrar latência
		collector.RecordHistogram("http_request_duration_seconds", latency, labels)

		// Registrar tamanho da resposta
		collector.RecordGauge("http_response_size_bytes", float64(size), labels)
	}
}
