package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"

	"github.com/devleo-m/go-zero/internal/domain/shared"
	"github.com/devleo-m/go-zero/internal/interface/http/dto"
)

// ErrorHandler handler centralizado de erros
type ErrorHandler struct {
	logger        *zap.Logger
	environment   string
	enableDetails bool
}

// ErrorHandlerConfig configuração do error handler
type ErrorHandlerConfig struct {
	Logger        *zap.Logger
	Environment   string
	EnableDetails bool
}

// NewErrorHandler cria uma nova instância do ErrorHandler
func NewErrorHandler(config ErrorHandlerConfig) *ErrorHandler {
	return &ErrorHandler{
		logger:        config.Logger,
		environment:   config.Environment,
		enableDetails: config.EnableDetails,
	}
}

// HandleError trata erros de forma centralizada
func (h *ErrorHandler) HandleError(c *gin.Context, err error) {
	// Log do erro
	h.logError(c, err)

	// Determinar tipo de erro e resposta apropriada
	switch {
	case isValidationError(err):
		h.handleValidationError(c, err)
	case isDomainError(err):
		h.handleDomainError(c, err)
	case isNotFoundError(err):
		h.handleNotFoundError(c, err)
	case isUnauthorizedError(err):
		h.handleUnauthorizedError(c, err)
	case isForbiddenError(err):
		h.handleForbiddenError(c, err)
	case isConflictError(err):
		h.handleConflictError(c, err)
	case isTimeoutError(err):
		h.handleTimeoutError(c, err)
	default:
		h.handleInternalError(c, err)
	}
}

// logError registra o erro no logger
func (h *ErrorHandler) logError(c *gin.Context, err error) {
	// Extrair informações do contexto
	userID, _ := c.Get("user_id")
	requestID, _ := c.Get("request_id")
	path := c.Request.URL.Path
	method := c.Request.Method

	fields := []zap.Field{
		zap.Error(err),
		zap.String("path", path),
		zap.String("method", method),
		zap.String("ip", c.ClientIP()),
		zap.String("user_agent", c.Request.UserAgent()),
	}

	if userID != nil {
		fields = append(fields, zap.String("user_id", userID.(string)))
	}

	if requestID != nil {
		fields = append(fields, zap.String("request_id", requestID.(string)))
	}

	// Determinar nível de log baseado no tipo de erro
	switch {
	case isValidationError(err):
		h.logger.Warn("Validation error", fields...)
	case isDomainError(err):
		h.logger.Warn("Domain error", fields...)
	case isNotFoundError(err):
		h.logger.Info("Not found error", fields...)
	case isUnauthorizedError(err):
		h.logger.Warn("Unauthorized error", fields...)
	case isForbiddenError(err):
		h.logger.Warn("Forbidden error", fields...)
	case isConflictError(err):
		h.logger.Warn("Conflict error", fields...)
	case isTimeoutError(err):
		h.logger.Warn("Timeout error", fields...)
	default:
		h.logger.Error("Internal error", fields...)
	}
}

// handleValidationError trata erros de validação
func (h *ErrorHandler) handleValidationError(c *gin.Context, err error) {
	var validationErrors []dto.ValidationErrorDetail

	if ve, ok := err.(validator.ValidationErrors); ok {
		for _, e := range ve {
			validationErrors = append(validationErrors, dto.ValidationErrorDetail{
				Field:   e.Field(),
				Tag:     e.Tag(),
				Value:   getFieldValue(e.Value()),
				Message: getValidationMessage(e),
			})
		}
	} else {
		// Erro de validação genérico
		validationErrors = append(validationErrors, dto.ValidationErrorDetail{
			Field:   "general",
			Message: err.Error(),
		})
	}

	response := dto.NewValidationErrorResponse(validationErrors)
	c.JSON(http.StatusBadRequest, response)
}

// handleDomainError trata erros de domínio
func (h *ErrorHandler) handleDomainError(c *gin.Context, err error) {
	var domainErr *shared.DomainError
	if errors.As(err, &domainErr) {
		response := dto.NewErrorResponse(
			domainErr.Code,
			domainErr.Message,
		)
		c.JSON(http.StatusBadRequest, response)
	} else {
		response := dto.NewErrorResponse(
			"DOMAIN_ERROR",
			err.Error(),
		)
		c.JSON(http.StatusBadRequest, response)
	}
}

// handleNotFoundError trata erros de não encontrado
func (h *ErrorHandler) handleNotFoundError(c *gin.Context, err error) {
	response := dto.NewErrorResponse(
		"NOT_FOUND",
		"Resource not found",
	)
	c.JSON(http.StatusNotFound, response)
}

// handleUnauthorizedError trata erros de não autorizado
func (h *ErrorHandler) handleUnauthorizedError(c *gin.Context, err error) {
	response := dto.NewErrorResponse(
		"UNAUTHORIZED",
		"Authentication required",
	)
	c.JSON(http.StatusUnauthorized, response)
}

// handleForbiddenError trata erros de proibido
func (h *ErrorHandler) handleForbiddenError(c *gin.Context, err error) {
	response := dto.NewErrorResponse(
		"FORBIDDEN",
		"Insufficient permissions",
	)
	c.JSON(http.StatusForbidden, response)
}

// handleConflictError trata erros de conflito
func (h *ErrorHandler) handleConflictError(c *gin.Context, err error) {
	response := dto.NewErrorResponse(
		"CONFLICT",
		"Resource already exists or conflict occurred",
	)
	c.JSON(http.StatusConflict, response)
}

// handleTimeoutError trata erros de timeout
func (h *ErrorHandler) handleTimeoutError(c *gin.Context, err error) {
	response := dto.NewErrorResponse(
		"TIMEOUT",
		"Request timeout",
	)
	c.JSON(http.StatusRequestTimeout, response)
}

// handleInternalError trata erros internos
func (h *ErrorHandler) handleInternalError(c *gin.Context, err error) {
	response := dto.NewErrorResponse(
		"INTERNAL_SERVER_ERROR",
		"An internal server error occurred",
	)
	c.JSON(http.StatusInternalServerError, response)
}

// ==========================================
// ERROR TYPE DETECTION
// ==========================================

// isValidationError verifica se é um erro de validação
func isValidationError(err error) bool {
	var ve validator.ValidationErrors
	return errors.As(err, &ve)
}

// isDomainError verifica se é um erro de domínio
func isDomainError(err error) bool {
	var domainErr *shared.DomainError
	return errors.As(err, &domainErr)
}

// isNotFoundError verifica se é um erro de não encontrado
func isNotFoundError(err error) bool {
	return strings.Contains(strings.ToLower(err.Error()), "not found") ||
		strings.Contains(strings.ToLower(err.Error()), "not exist")
}

// isUnauthorizedError verifica se é um erro de não autorizado
func isUnauthorizedError(err error) bool {
	return strings.Contains(strings.ToLower(err.Error()), "unauthorized") ||
		strings.Contains(strings.ToLower(err.Error()), "invalid credentials") ||
		strings.Contains(strings.ToLower(err.Error()), "authentication")
}

// isForbiddenError verifica se é um erro de proibido
func isForbiddenError(err error) bool {
	return strings.Contains(strings.ToLower(err.Error()), "forbidden") ||
		strings.Contains(strings.ToLower(err.Error()), "permission") ||
		strings.Contains(strings.ToLower(err.Error()), "access denied")
}

// isConflictError verifica se é um erro de conflito
func isConflictError(err error) bool {
	return strings.Contains(strings.ToLower(err.Error()), "conflict") ||
		strings.Contains(strings.ToLower(err.Error()), "already exists") ||
		strings.Contains(strings.ToLower(err.Error()), "duplicate")
}

// isTimeoutError verifica se é um erro de timeout
func isTimeoutError(err error) bool {
	return strings.Contains(strings.ToLower(err.Error()), "timeout") ||
		strings.Contains(strings.ToLower(err.Error()), "deadline exceeded")
}

// ==========================================
// HELPER FUNCTIONS
// ==========================================

// getFieldValue obtém o valor do campo para exibição
func getFieldValue(value interface{}) string {
	if value == nil {
		return ""
	}
	return strings.TrimSpace(strings.Trim(strings.Trim(fmt.Sprintf("%v", value), "[]"), "{}"))
}

// getValidationMessage obtém mensagem de validação personalizada
func getValidationMessage(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return e.Field() + " is required"
	case "email":
		return e.Field() + " must be a valid email address"
	case "min":
		return e.Field() + " must be at least " + e.Param() + " characters long"
	case "max":
		return e.Field() + " must be at most " + e.Param() + " characters long"
	case "oneof":
		return e.Field() + " must be one of: " + e.Param()
	case "uuid4":
		return e.Field() + " must be a valid UUID"
	case "len":
		return e.Field() + " must be exactly " + e.Param() + " characters long"
	case "gte":
		return e.Field() + " must be greater than or equal to " + e.Param()
	case "lte":
		return e.Field() + " must be less than or equal to " + e.Param()
	case "gt":
		return e.Field() + " must be greater than " + e.Param()
	case "lt":
		return e.Field() + " must be less than " + e.Param()
	case "alphanum":
		return e.Field() + " must contain only alphanumeric characters"
	case "alpha":
		return e.Field() + " must contain only alphabetic characters"
	case "numeric":
		return e.Field() + " must be numeric"
	case "url":
		return e.Field() + " must be a valid URL"
	case "datetime":
		return e.Field() + " must be a valid datetime"
	case "date":
		return e.Field() + " must be a valid date"
	case "time":
		return e.Field() + " must be a valid time"
	default:
		return e.Field() + " is invalid"
	}
}

// ==========================================
// GIN ERROR HANDLER MIDDLEWARE
// ==========================================

// ErrorHandlerMiddleware middleware para capturar erros do Gin
func (h *ErrorHandler) ErrorHandlerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Verificar se há erros
		if len(c.Errors) > 0 {
			err := c.Errors.Last()
			h.HandleError(c, err.Err)
			c.Abort()
		}
	}
}

// ==========================================
// CUSTOM ERROR RESPONSES
// ==========================================

// HandleCustomError trata erros customizados com códigos específicos
func (h *ErrorHandler) HandleCustomError(c *gin.Context, code string, message string, statusCode int) {
	response := dto.NewErrorResponse(code, message)
	c.JSON(statusCode, response)
}

// HandleValidationError trata erros de validação específicos
func (h *ErrorHandler) HandleValidationError(c *gin.Context, field string, message string) {
	validationErrors := []dto.ValidationErrorDetail{
		{
			Field:   field,
			Message: message,
		},
	}

	response := dto.NewValidationErrorResponse(validationErrors)
	c.JSON(http.StatusBadRequest, response)
}

// HandleBusinessError trata erros de negócio específicos
func (h *ErrorHandler) HandleBusinessError(c *gin.Context, code string, message string) {
	response := dto.NewErrorResponse(code, message)
	c.JSON(http.StatusUnprocessableEntity, response)
}

// ==========================================
// ENHANCED ERROR HANDLING
// ==========================================

// ErrorResponseWithDetails representa uma resposta de erro com detalhes
type ErrorResponseWithDetails struct {
	Success   bool        `json:"success" example:"false"`
	Error     string      `json:"error" example:"VALIDATION_ERROR"`
	Message   string      `json:"message" example:"Validation failed"`
	Details   interface{} `json:"details,omitempty"`
	RequestID string      `json:"request_id,omitempty"`
	Timestamp string      `json:"timestamp"`
	TraceID   string      `json:"trace_id,omitempty"`
}

// HandleErrorWithDetails trata erros com detalhes adicionais
func (h *ErrorHandler) HandleErrorWithDetails(c *gin.Context, err error, details interface{}) {
	// Log do erro
	h.logError(c, err)

	// Extrair informações do contexto
	requestID, _ := c.Get("request_id")
	traceID, _ := c.Get("trace_id")

	response := ErrorResponseWithDetails{
		Success:   false,
		Error:     h.getErrorCode(err),
		Message:   h.getErrorMessage(err),
		Details:   details,
		RequestID: getStringValue(requestID),
		Timestamp: time.Now().Format(time.RFC3339),
		TraceID:   getStringValue(traceID),
	}

	// Determinar status code
	statusCode := h.getStatusCode(err)
	c.JSON(statusCode, response)
}

// HandleRetryableError trata erros que podem ser retentados
func (h *ErrorHandler) HandleRetryableError(c *gin.Context, err error, retryAfter int) {
	// Log do erro
	h.logError(c, err)

	response := dto.NewErrorResponse(
		"RETRYABLE_ERROR",
		"Service temporarily unavailable, please retry later",
	)

	// Adicionar header de retry
	c.Header("Retry-After", fmt.Sprintf("%d", retryAfter))
	c.JSON(http.StatusServiceUnavailable, response)
}

// HandleCircuitBreakerError trata erros de circuit breaker
func (h *ErrorHandler) HandleCircuitBreakerError(c *gin.Context, service string) {
	response := dto.NewErrorResponse(
		"CIRCUIT_BREAKER_OPEN",
		fmt.Sprintf("Service %s is temporarily unavailable due to high error rate", service),
	)

	c.JSON(http.StatusServiceUnavailable, response)
}

// ==========================================
// HELPER FUNCTIONS
// ==========================================

// getErrorCode extrai o código do erro
func (h *ErrorHandler) getErrorCode(err error) string {
	if domainErr, ok := err.(*shared.DomainError); ok {
		return domainErr.Code
	}

	// Mapear tipos de erro para códigos
	switch {
	case isValidationError(err):
		return "VALIDATION_ERROR"
	case isNotFoundError(err):
		return "NOT_FOUND"
	case isUnauthorizedError(err):
		return "UNAUTHORIZED"
	case isForbiddenError(err):
		return "FORBIDDEN"
	case isConflictError(err):
		return "CONFLICT"
	case isTimeoutError(err):
		return "TIMEOUT"
	default:
		return "INTERNAL_ERROR"
	}
}

// getErrorMessage extrai a mensagem do erro
func (h *ErrorHandler) getErrorMessage(err error) string {
	if domainErr, ok := err.(*shared.DomainError); ok {
		return domainErr.Message
	}

	// Mapear tipos de erro para mensagens
	switch {
	case isValidationError(err):
		return "Validation failed"
	case isNotFoundError(err):
		return "Resource not found"
	case isUnauthorizedError(err):
		return "Authentication required"
	case isForbiddenError(err):
		return "Insufficient permissions"
	case isConflictError(err):
		return "Resource conflict"
	case isTimeoutError(err):
		return "Request timeout"
	default:
		return "An internal error occurred"
	}
}

// getStatusCode determina o status HTTP baseado no erro
func (h *ErrorHandler) getStatusCode(err error) int {
	switch {
	case isValidationError(err):
		return http.StatusBadRequest
	case isNotFoundError(err):
		return http.StatusNotFound
	case isUnauthorizedError(err):
		return http.StatusUnauthorized
	case isForbiddenError(err):
		return http.StatusForbidden
	case isConflictError(err):
		return http.StatusConflict
	case isTimeoutError(err):
		return http.StatusRequestTimeout
	default:
		return http.StatusInternalServerError
	}
}

// getStringValue converte interface{} para string de forma segura
func getStringValue(value interface{}) string {
	if value == nil {
		return ""
	}
	if str, ok := value.(string); ok {
		return str
	}
	return fmt.Sprintf("%v", value)
}
