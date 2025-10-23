package middleware

import (
	"bytes"
	"io"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// responseWriter wrapper para capturar o status code
type responseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

// Write sobrescreve o método Write para capturar o body da resposta
func (w responseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// WriteString sobrescreve o método WriteString
func (w responseWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

// LoggerMiddleware middleware de logging para requisições HTTP
func LoggerMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Capturar o body da requisição
		var requestBody []byte
		if c.Request.Body != nil {
			requestBody, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		// Criar wrapper para capturar a resposta
		blw := &responseWriter{
			ResponseWriter: c.Writer,
			body:           &bytes.Buffer{},
		}
		c.Writer = blw

		// Processar a requisição
		c.Next()

		// Calcular duração
		latency := time.Since(start)

		// Extrair informações do contexto
		userID, _ := c.Get("user_id")
		userRole, _ := c.Get("user_role")
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = c.GetString("request_id")
		}

		// Construir query string completa
		if raw != "" {
			path = path + "?" + raw
		}

		// Log da requisição
		fields := []zap.Field{
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.Int("status", c.Writer.Status()),
			zap.Duration("latency", latency),
			zap.String("ip", c.ClientIP()),
			zap.String("user_agent", c.Request.UserAgent()),
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
		}

		// Adicionar campos opcionais
		if requestID != "" {
			fields = append(fields, zap.String("request_id", requestID))
		}

		if userID != nil {
			fields = append(fields, zap.String("user_id", userID.(string)))
		}

		if userRole != nil {
			fields = append(fields, zap.String("user_role", userRole.(string)))
		}

		// Adicionar informações da requisição em debug
		if logger.Core().Enabled(zap.DebugLevel) {
			fields = append(fields,
				zap.String("request_body", string(requestBody)),
				zap.String("response_body", blw.body.String()),
				zap.String("content_type", c.GetHeader("Content-Type")),
				zap.String("referer", c.GetHeader("Referer")),
			)
		}

		// Determinar nível de log baseado no status code
		switch {
		case c.Writer.Status() >= 500:
			logger.Error("HTTP Request", fields...)
		case c.Writer.Status() >= 400:
			logger.Warn("HTTP Request", fields...)
		default:
			logger.Info("HTTP Request", fields...)
		}
	}
}

// SecurityHeadersMiddleware middleware para adicionar headers de segurança
func SecurityHeadersMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Prevenir clickjacking
		c.Header("X-Frame-Options", "DENY")

		// Prevenir MIME type sniffing
		c.Header("X-Content-Type-Options", "nosniff")

		// Habilitar XSS protection
		c.Header("X-XSS-Protection", "1; mode=block")

		// Referrer policy
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")

		// Content Security Policy (básico)
		c.Header("Content-Security-Policy", "default-src 'self'")

		// Strict Transport Security (apenas em HTTPS)
		if c.Request.TLS != nil {
			c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		}

		c.Next()
	}
}

// LogBusinessEvent loga eventos de negócio
func LogBusinessEvent(logger *zap.Logger, event string, userID string, data map[string]interface{}) {
	fields := []zap.Field{
		zap.String("event", event),
		zap.String("user_id", userID),
		zap.String("timestamp", time.Now().Format(time.RFC3339)),
	}

	// Adicionar dados do evento
	for key, value := range data {
		fields = append(fields, zap.Any(key, value))
	}

	logger.Info("Business Event", fields...)
}

// LogSecurityEvent loga eventos de segurança
func LogSecurityEvent(logger *zap.Logger, event string, userID string, ip string, userAgent string, err error) {
	fields := []zap.Field{
		zap.String("event", event),
		zap.String("user_id", userID),
		zap.String("ip", ip),
		zap.String("user_agent", userAgent),
		zap.String("timestamp", time.Now().Format(time.RFC3339)),
	}

	if err != nil {
		fields = append(fields, zap.Error(err))
		logger.Warn("Security Event", fields...)
	} else {
		logger.Info("Security Event", fields...)
	}
}
