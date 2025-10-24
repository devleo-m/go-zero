package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// RateLimiter representa um limitador de taxa
type RateLimiter struct {
	requests map[string][]time.Time
	mutex    sync.RWMutex
	limit    int
	window   time.Duration
}

// NewRateLimiter cria um novo limitador de taxa
func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		requests: make(map[string][]time.Time),
		limit:    limit,
		window:   window,
	}
}

// RateLimit cria um middleware de rate limiting
func RateLimit(limiter *RateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Obter identificador do cliente (IP ou user ID)
		clientID := getClientIdentifier(c)

		// Verificar se o cliente excedeu o limite
		if !limiter.Allow(clientID) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"success": false,
				"error":   "RATE_LIMIT_EXCEEDED",
				"message": "Too many requests",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// Allow verifica se uma requisição é permitida
func (rl *RateLimiter) Allow(clientID string) bool {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	now := time.Now()

	// Limpar requisições antigas
	rl.cleanup(clientID, now)

	// Verificar se ainda há espaço para mais requisições
	if len(rl.requests[clientID]) >= rl.limit {
		return false
	}

	// Adicionar nova requisição
	rl.requests[clientID] = append(rl.requests[clientID], now)
	return true
}

// cleanup remove requisições antigas
func (rl *RateLimiter) cleanup(clientID string, now time.Time) {
	cutoff := now.Add(-rl.window)
	requests := rl.requests[clientID]

	// Encontrar o primeiro índice que não deve ser removido
	start := 0
	for i, reqTime := range requests {
		if reqTime.After(cutoff) {
			start = i
			break
		}
	}

	// Manter apenas as requisições dentro da janela
	rl.requests[clientID] = requests[start:]
}

// getClientIdentifier obtém um identificador único para o cliente
func getClientIdentifier(c *gin.Context) string {
	// Tentar obter user ID se estiver autenticado
	if userID, exists := c.Get("user_id"); exists {
		if id, ok := userID.(string); ok {
			return "user:" + id
		}
	}

	// Usar IP como fallback
	return "ip:" + c.ClientIP()
}

// GetRemainingRequests retorna o número de requisições restantes
func (rl *RateLimiter) GetRemainingRequests(clientID string) int {
	rl.mutex.RLock()
	defer rl.mutex.RUnlock()

	now := time.Now()
	cutoff := now.Add(-rl.window)
	requests := rl.requests[clientID]

	// Contar requisições dentro da janela
	count := 0
	for _, reqTime := range requests {
		if reqTime.After(cutoff) {
			count++
		}
	}

	return rl.limit - count
}

// GetResetTime retorna o tempo até o reset do rate limit
func (rl *RateLimiter) GetResetTime(clientID string) time.Time {
	rl.mutex.RLock()
	defer rl.mutex.RUnlock()

	requests := rl.requests[clientID]
	if len(requests) == 0 {
		return time.Now()
	}

	// Retornar o tempo da requisição mais antiga + janela
	oldest := requests[0]
	return oldest.Add(rl.window)
}
