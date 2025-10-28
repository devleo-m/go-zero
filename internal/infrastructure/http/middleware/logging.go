package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// LoggingMiddleware cria um middleware de logging.
func LoggingMiddleware(logger interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Gerar request ID único
		requestID := uuid.New().String()
		c.Set("request_id", requestID)

		// Adicionar request ID ao header de resposta
		c.Header("X-Request-ID", requestID)

		// Registrar início da requisição
		start := time.Now()

		// Processar requisição
		c.Next()

		// Calculator duração
		duration := time.Since(start)

		// Registrar fim da requisição
		logRequest(c, start, duration)
	}
}

// logRequest registra informações da requisição.
func logRequest(c *gin.Context, start time.Time, duration time.Duration) {
	// Obter informações da requisição
	method := c.Request.Method
	path := c.Request.URL.Path
	statusCode := c.Writer.Status()
	clientIP := c.ClientIP()
	userAgent := c.Request.UserAgent()

	// Obter request ID
	requestID, _ := c.Get("request_id")

	// Obter informações do usuário se autenticado
	userRole, _ := c.Get("user_role")

	// Criar campos de log
	fields := []interface{}{
		"method", method,
		"path", path,
		"status", statusCode,
		"duration", duration.String(),
		"client_ip", clientIP,
		"user_agent", userAgent,
		"request_id", requestID,
	}

	// Adicionar informações do usuário se disponíveis
	if userRole != nil {
		fields = append(fields, "user_role", userRole)
	}

	// Adicionar tamanho da resposta
	fields = append(fields, "response_size", c.Writer.Size())

	// Log baseado no status code
	if statusCode >= 500 {
		// Erro do servidor
		logError(c, "HTTP Request", fields...)
	} else if statusCode >= 400 {
		// Erro do cliente
		logWarn(c, "HTTP Request", fields...)
	} else {
		// Sucesso
		logInfo(c, "HTTP Request", fields...)
	}
}

// logInfo registra um log de informação.
func logInfo(c *gin.Context, msg string, fields ...interface{}) {
	// Aqui você pode integrar com seu sistema de logging
	// Por exemplo, usando zap, logrus, etc.
	// Por enquanto, vamos usar o logger padrão do Gin
	gin.DefaultWriter.Write([]byte(msg + "\n"))
}

// logWarn registra um log de aviso.
func logWarn(c *gin.Context, msg string, fields ...interface{}) {
	gin.DefaultErrorWriter.Write([]byte(msg + "\n"))
}

// logError registra um log de erro.
func logError(c *gin.Context, msg string, fields ...interface{}) {
	gin.DefaultErrorWriter.Write([]byte(msg + "\n"))
}

// RequestIDMiddleware adiciona um request ID único a cada requisição.
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}

		c.Set("request_id", requestID)
		c.Header("X-Request-ID", requestID)
		c.Next()
	}
}

// RecoveryMiddleware cria um middleware de recovery personalizado.
func RecoveryMiddleware() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		// Log do erro
		requestID, _ := c.Get("request_id")

		// Aqui você pode integrar com seu sistema de logging
		gin.DefaultErrorWriter.Write([]byte("Panic recovered: " + recovered.(error).Error() + "\n"))

		// Responder com erro interno do servidor
		c.JSON(500, gin.H{
			"success":    false,
			"error":      "INTERNAL_SERVER_ERROR",
			"message":    "An internal error occurred",
			"request_id": requestID,
		})
	})
}
