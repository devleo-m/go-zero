package middleware

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/devleo-m/go-zero/internal/interface/http/dto"
)

// RecoveryMiddleware middleware de recovery para capturar panics
func RecoveryMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Verificar se é uma conexão quebrada
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") ||
							strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				// Obter informações da requisição
				httpRequest, _ := httputil.DumpRequest(c.Request, false)

				// Extrair informações do contexto
				userID, _ := c.Get("user_id")
				requestID, _ := c.Get("request_id")

				// Log do panic
				fields := []zap.Field{
					zap.Any("error", err),
					zap.String("request", string(httpRequest)),
					zap.String("stack", string(debug.Stack())),
					zap.String("timestamp", time.Now().Format(time.RFC3339)),
				}

				if userID != nil {
					fields = append(fields, zap.String("user_id", userID.(string)))
				}

				if requestID != nil {
					fields = append(fields, zap.String("request_id", requestID.(string)))
				}

				// Log baseado no tipo de erro
				if brokenPipe {
					logger.Error("Broken pipe error", fields...)
				} else {
					logger.Error("Panic recovered", fields...)
				}

				// Se a conexão foi quebrada, não enviar resposta
				if brokenPipe {
					c.Error(err.(error))
					c.Abort()
					return
				}

				// Enviar resposta de erro
				errorResponse := dto.NewErrorResponse(
					"INTERNAL_SERVER_ERROR",
					"An internal server error occurred",
				)

				c.JSON(http.StatusInternalServerError, errorResponse)
				c.Abort()
			}
		}()

		c.Next()
	}
}

// CustomRecoveryMiddleware middleware de recovery customizado com mais opções
func CustomRecoveryMiddleware(logger *zap.Logger, customRecovery func(c *gin.Context, recovered interface{})) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Verificar se é uma conexão quebrada
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") ||
							strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				// Se a conexão foi quebrada, não enviar resposta
				if brokenPipe {
					c.Error(err.(error))
					c.Abort()
					return
				}

				// Chamar função de recovery customizada
				if customRecovery != nil {
					customRecovery(c, err)
				} else {
					// Recovery padrão
					errorResponse := dto.NewErrorResponse(
						"INTERNAL_SERVER_ERROR",
						"An internal server error occurred",
					)
					c.JSON(http.StatusInternalServerError, errorResponse)
				}

				c.Abort()
			}
		}()

		c.Next()
	}
}

// DefaultRecoveryHandler handler de recovery padrão
func DefaultRecoveryHandler(logger *zap.Logger) func(c *gin.Context, recovered interface{}) {
	return func(c *gin.Context, recovered interface{}) {
		// Obter informações da requisição
		httpRequest, _ := httputil.DumpRequest(c.Request, false)

		// Extrair informações do contexto
		userID, _ := c.Get("user_id")
		requestID, _ := c.Get("request_id")

		// Log do panic
		fields := []zap.Field{
			zap.Any("error", recovered),
			zap.String("request", string(httpRequest)),
			zap.String("stack", string(debug.Stack())),
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
		}

		if userID != nil {
			fields = append(fields, zap.String("user_id", userID.(string)))
		}

		if requestID != nil {
			fields = append(fields, zap.String("request_id", requestID.(string)))
		}

		logger.Error("Panic recovered", fields...)

		// Enviar resposta de erro
		errorResponse := dto.NewErrorResponse(
			"INTERNAL_SERVER_ERROR",
			"An internal server error occurred",
		)

		c.JSON(http.StatusInternalServerError, errorResponse)
	}
}

// DevelopmentRecoveryHandler handler de recovery para desenvolvimento
func DevelopmentRecoveryHandler(logger *zap.Logger) func(c *gin.Context, recovered interface{}) {
	return func(c *gin.Context, recovered interface{}) {
		// Obter informações da requisição
		httpRequest, _ := httputil.DumpRequest(c.Request, false)

		// Extrair informações do contexto
		userID, _ := c.Get("user_id")
		requestID, _ := c.Get("request_id")

		// Log detalhado do panic
		fields := []zap.Field{
			zap.Any("error", recovered),
			zap.String("request", string(httpRequest)),
			zap.String("stack", string(debug.Stack())),
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
		}

		if userID != nil {
			fields = append(fields, zap.String("user_id", userID.(string)))
		}

		if requestID != nil {
			fields = append(fields, zap.String("request_id", requestID.(string)))
		}

		logger.Error("Panic recovered (Development)", fields...)

		// Resposta detalhada para desenvolvimento
		errorResponse := dto.ErrorResponse{
			Success: false,
			Error:   "INTERNAL_SERVER_ERROR",
			Message: fmt.Sprintf("Panic: %v", recovered),
		}

		c.JSON(http.StatusInternalServerError, errorResponse)
	}
}

// ProductionRecoveryHandler handler de recovery para produção
func ProductionRecoveryHandler(logger *zap.Logger) func(c *gin.Context, recovered interface{}) {
	return func(c *gin.Context, recovered interface{}) {
		// Obter informações da requisição
		httpRequest, _ := httputil.DumpRequest(c.Request, false)

		// Extrair informações do contexto
		userID, _ := c.Get("user_id")
		requestID, _ := c.Get("request_id")

		// Log do panic (sem stack trace em produção)
		fields := []zap.Field{
			zap.Any("error", recovered),
			zap.String("request", string(httpRequest)),
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
		}

		if userID != nil {
			fields = append(fields, zap.String("user_id", userID.(string)))
		}

		if requestID != nil {
			fields = append(fields, zap.String("request_id", requestID.(string)))
		}

		logger.Error("Panic recovered (Production)", fields...)

		// Resposta genérica para produção
		errorResponse := dto.NewErrorResponse(
			"INTERNAL_SERVER_ERROR",
			"An internal server error occurred. Please try again later.",
		)

		c.JSON(http.StatusInternalServerError, errorResponse)
	}
}

// TimeoutMiddleware middleware para timeout de requisições
func TimeoutMiddleware(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Criar contexto com timeout
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()

		// Substituir o contexto da requisição
		c.Request = c.Request.WithContext(ctx)

		// Canal para sinalizar conclusão
		done := make(chan bool, 1)

		go func() {
			c.Next()
			done <- true
		}()

		select {
		case <-done:
			// Requisição concluída normalmente
		case <-ctx.Done():
			// Timeout atingido
			c.JSON(http.StatusRequestTimeout, dto.NewErrorResponse(
				"REQUEST_TIMEOUT",
				"Request timeout. Please try again.",
			))
			c.Abort()
		}
	}
}

// MaxBodySizeMiddleware middleware para limitar tamanho do body da requisição
func MaxBodySizeMiddleware(maxSize int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Verificar Content-Length
		if c.Request.ContentLength > maxSize {
			c.JSON(http.StatusRequestEntityTooLarge, dto.NewErrorResponse(
				"REQUEST_TOO_LARGE",
				"Request body too large",
			))
			c.Abort()
			return
		}

		// Limitar o body reader
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxSize)

		c.Next()
	}
}
