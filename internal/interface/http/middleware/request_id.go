package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// RequestIDMiddleware middleware para adicionar Request ID a todas as requisições
func RequestIDMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Tentar obter Request ID do header
		requestID := c.GetHeader("X-Request-ID")

		// Se não existir, gerar um novo
		if requestID == "" {
			requestID = generateRequestID()
		}

		// Adicionar ao contexto
		c.Set("request_id", requestID)

		// Adicionar ao header de resposta
		c.Header("X-Request-ID", requestID)

		// Log da requisição
		logger.Info("Request started",
			zap.String("request_id", requestID),
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.String("ip", c.ClientIP()),
			zap.String("user_agent", c.Request.UserAgent()),
		)

		c.Next()

		// Log da resposta
		logger.Info("Request completed",
			zap.String("request_id", requestID),
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.Int("status", c.Writer.Status()),
			zap.Int("size", c.Writer.Size()),
		)
	}
}

// generateRequestID gera um ID único para a requisição
func generateRequestID() string {
	// Gerar 16 bytes aleatórios
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		// Fallback para timestamp se falhar
		return "req_" + hex.EncodeToString([]byte(time.Now().Format("20060102150405")))
	}

	// Converter para hex e adicionar prefixo
	return "req_" + hex.EncodeToString(bytes)
}
