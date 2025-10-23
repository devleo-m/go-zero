package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/devleo-m/go-zero/internal/interface/http/dto"
)

// RateLimiter interface para rate limiting
type RateLimiter interface {
	Allow(ctx context.Context, key string, limit int, window time.Duration) (bool, error)
	GetRemaining(ctx context.Context, key string, limit int, window time.Duration) (int, error)
	GetResetTime(ctx context.Context, key string, window time.Duration) (time.Time, error)
}

// RateLimitConfig configuração do rate limiting
type RateLimitConfig struct {
	DefaultLimit  int           // Limite padrão por minuto
	DefaultWindow time.Duration // Janela de tempo padrão
	SkipPaths     []string      // Caminhos que devem ser ignorados
	SkipIPs       []string      // IPs que devem ser ignorados
}

// RateLimitMiddleware middleware para rate limiting
func RateLimitMiddleware(limiter RateLimiter, config RateLimitConfig, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Verificar se deve pular o rate limiting
		if shouldSkipRateLimit(c, config) {
			c.Next()
			return
		}

		// Obter chave para rate limiting
		key := getRateLimitKey(c)

		// Verificar limite
		allowed, err := limiter.Allow(c.Request.Context(), key, config.DefaultLimit, config.DefaultWindow)
		if err != nil {
			logger.Error("Rate limiter error",
				zap.String("key", key),
				zap.Error(err),
			)
			// Em caso de erro, permitir a requisição
			c.Next()
			return
		}

		if !allowed {
			// Obter informações de rate limit
			remaining, _ := limiter.GetRemaining(c.Request.Context(), key, config.DefaultLimit, config.DefaultWindow)
			resetTime, _ := limiter.GetResetTime(c.Request.Context(), key, config.DefaultWindow)

			// Adicionar headers de rate limit
			c.Header("X-RateLimit-Limit", strconv.Itoa(config.DefaultLimit))
			c.Header("X-RateLimit-Remaining", strconv.Itoa(remaining))
			c.Header("X-RateLimit-Reset", strconv.FormatInt(resetTime.Unix(), 10))
			c.Header("Retry-After", strconv.FormatInt(int64(time.Until(resetTime).Seconds()), 10))

			// Log do rate limit
			logger.Warn("Rate limit exceeded",
				zap.String("key", key),
				zap.String("ip", c.ClientIP()),
				zap.String("path", c.Request.URL.Path),
				zap.Int("limit", config.DefaultLimit),
				zap.Int("remaining", remaining),
			)

			// Retornar erro 429
			response := dto.NewErrorResponse(
				"RATE_LIMIT_EXCEEDED",
				fmt.Sprintf("Rate limit exceeded. Try again in %d seconds", int(time.Until(resetTime).Seconds())),
			)
			c.JSON(http.StatusTooManyRequests, response)
			c.Abort()
			return
		}

		// Obter informações de rate limit para headers
		remaining, _ := limiter.GetRemaining(c.Request.Context(), key, config.DefaultLimit, config.DefaultWindow)
		resetTime, _ := limiter.GetResetTime(c.Request.Context(), key, config.DefaultWindow)

		// Adicionar headers de rate limit
		c.Header("X-RateLimit-Limit", strconv.Itoa(config.DefaultLimit))
		c.Header("X-RateLimit-Remaining", strconv.Itoa(remaining))
		c.Header("X-RateLimit-Reset", strconv.FormatInt(resetTime.Unix(), 10))

		c.Next()
	}
}

// shouldSkipRateLimit verifica se deve pular o rate limiting
func shouldSkipRateLimit(c *gin.Context, config RateLimitConfig) bool {
	path := c.Request.URL.Path
	ip := c.ClientIP()

	// Verificar caminhos ignorados
	for _, skipPath := range config.SkipPaths {
		if path == skipPath {
			return true
		}
	}

	// Verificar IPs ignorados
	for _, skipIP := range config.SkipIPs {
		if ip == skipIP {
			return true
		}
	}

	return false
}

// getRateLimitKey gera a chave para rate limiting
func getRateLimitKey(c *gin.Context) string {
	ip := c.ClientIP()
	userID, exists := c.Get("user_id")

	if exists {
		// Rate limit por usuário autenticado
		return fmt.Sprintf("user:%s", userID)
	}

	// Rate limit por IP
	return fmt.Sprintf("ip:%s", ip)
}

// ==========================================
// RATE LIMIT BY ENDPOINT
// ==========================================

// EndpointRateLimitConfig configuração de rate limit por endpoint
type EndpointRateLimitConfig struct {
	Path   string
	Limit  int
	Window time.Duration
}

// RateLimitByEndpointMiddleware middleware para rate limiting por endpoint
func RateLimitByEndpointMiddleware(limiter RateLimiter, configs map[string]EndpointRateLimitConfig, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path

		// Verificar se há configuração específica para este endpoint
		config, exists := configs[path]
		if !exists {
			c.Next()
			return
		}

		// Obter chave para rate limiting
		key := getRateLimitKey(c) + ":" + path

		// Verificar limite
		allowed, err := limiter.Allow(c.Request.Context(), key, config.Limit, config.Window)
		if err != nil {
			logger.Error("Rate limiter error",
				zap.String("key", key),
				zap.String("path", path),
				zap.Error(err),
			)
			c.Next()
			return
		}

		if !allowed {
			// Obter informações de rate limit
			remaining, _ := limiter.GetRemaining(c.Request.Context(), key, config.Limit, config.Window)
			resetTime, _ := limiter.GetResetTime(c.Request.Context(), key, config.Window)

			// Adicionar headers de rate limit
			c.Header("X-RateLimit-Limit", strconv.Itoa(config.Limit))
			c.Header("X-RateLimit-Remaining", strconv.Itoa(remaining))
			c.Header("X-RateLimit-Reset", strconv.FormatInt(resetTime.Unix(), 10))
			c.Header("Retry-After", strconv.FormatInt(int64(time.Until(resetTime).Seconds()), 10))

			// Log do rate limit
			logger.Warn("Rate limit exceeded for endpoint",
				zap.String("key", key),
				zap.String("path", path),
				zap.String("ip", c.ClientIP()),
				zap.Int("limit", config.Limit),
				zap.Int("remaining", remaining),
			)

			// Retornar erro 429
			response := dto.NewErrorResponse(
				"RATE_LIMIT_EXCEEDED",
				fmt.Sprintf("Rate limit exceeded for endpoint %s. Try again in %d seconds", path, int(time.Until(resetTime).Seconds())),
			)
			c.JSON(http.StatusTooManyRequests, response)
			c.Abort()
			return
		}

		// Obter informações de rate limit para headers
		remaining, _ := limiter.GetRemaining(c.Request.Context(), key, config.Limit, config.Window)
		resetTime, _ := limiter.GetResetTime(c.Request.Context(), key, config.Window)

		// Adicionar headers de rate limit
		c.Header("X-RateLimit-Limit", strconv.Itoa(config.Limit))
		c.Header("X-RateLimit-Remaining", strconv.Itoa(remaining))
		c.Header("X-RateLimit-Reset", strconv.FormatInt(resetTime.Unix(), 10))

		c.Next()
	}
}
