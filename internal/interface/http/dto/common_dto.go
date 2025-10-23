package dto

import (
	"time"
)

// ==========================================
// COMMON RESPONSE DTOs
// ==========================================

// SuccessResponse representa uma resposta de sucesso genérica
type SuccessResponse struct {
	Success   bool        `json:"success" example:"true"`
	Message   string      `json:"message" example:"Operation completed successfully"`
	Data      interface{} `json:"data,omitempty"`
	RequestID string      `json:"request_id,omitempty" example:"req_abc123def456"`
	Timestamp time.Time   `json:"timestamp" example:"2024-01-15T10:30:00Z"`
}

// ErrorResponse representa uma resposta de erro genérica
type ErrorResponse struct {
	Success   bool      `json:"success" example:"false"`
	Error     string    `json:"error" example:"VALIDATION_ERROR"`
	Message   string    `json:"message" example:"Validation failed"`
	RequestID string    `json:"request_id,omitempty" example:"req_abc123def456"`
	Timestamp time.Time `json:"timestamp" example:"2024-01-15T10:30:00Z"`
	TraceID   string    `json:"trace_id,omitempty" example:"trace_xyz789"`
}

// HealthCheckResponse representa a resposta do health check
type HealthCheckResponse struct {
	Success   bool      `json:"success" example:"true"`
	Message   string    `json:"message" example:"Service is healthy"`
	Timestamp time.Time `json:"timestamp" example:"2024-01-15T10:30:00Z"`
	Version   string    `json:"version" example:"1.0.0"`
	Uptime    string    `json:"uptime" example:"2h30m15s"`
	RequestID string    `json:"request_id,omitempty" example:"req_abc123def456"`
}

// ==========================================
// VALIDATION HELPERS
// ==========================================

// ValidationErrorDetail representa detalhes de um erro de validação
type ValidationErrorDetail struct {
	Field   string `json:"field" example:"email"`
	Tag     string `json:"tag" example:"email"`
	Value   string `json:"value" example:"invalid-email"`
	Message string `json:"message" example:"Invalid email format"`
}

// ValidationErrorResponse representa uma resposta de erro de validação
type ValidationErrorResponse struct {
	Success bool                    `json:"success" example:"false"`
	Error   string                  `json:"error" example:"Validation failed"`
	Message string                  `json:"message" example:"One or more fields are invalid"`
	Details []ValidationErrorDetail `json:"details"`
}

// ==========================================
// PAGINATION HELPERS
// ==========================================

// PaginationParams representa parâmetros de paginação
type PaginationParams struct {
	Page  int `form:"page" binding:"min=1" example:"1"`
	Limit int `form:"limit" binding:"min=1,max=100" example:"10"`
}

// PaginationResponse representa uma resposta paginada
type PaginationResponse struct {
	Data       interface{} `json:"data"`
	Pagination Pagination  `json:"pagination"`
}

// Pagination representa informações de paginação
type Pagination struct {
	Page       int   `json:"page" example:"1"`
	Limit      int   `json:"limit" example:"10"`
	Total      int64 `json:"total" example:"100"`
	TotalPages int   `json:"total_pages" example:"10"`
	HasNext    bool  `json:"has_next" example:"true"`
	HasPrev    bool  `json:"has_prev" example:"false"`
}

// GetPaginationDefaults retorna valores padrão para paginação
func GetPaginationDefaults() (page, limit int) {
	return 1, 10
}

// ValidatePagination valida e ajusta parâmetros de paginação
func ValidatePagination(page, limit int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	return page, limit
}

// NewPaginationResponse cria uma resposta paginada
func NewPaginationResponse(data interface{}, page, limit int, total int64) PaginationResponse {
	totalPages := int((total + int64(limit) - 1) / int64(limit))

	return PaginationResponse{
		Data: data,
		Pagination: Pagination{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: totalPages,
			HasNext:    page < totalPages,
			HasPrev:    page > 1,
		},
	}
}

// ==========================================
// SEARCH HELPERS
// ==========================================

// SearchParams representa parâmetros de busca
type SearchParams struct {
	Search string `form:"search" binding:"omitempty,max=100" example:"joão"`
	Sort   string `form:"sort" binding:"omitempty,oneof=name email created_at updated_at" example:"created_at"`
	Order  string `form:"order" binding:"omitempty,oneof=asc desc" example:"desc"`
}

// GetSearchDefaults retorna valores padrão para busca
func GetSearchDefaults() (search, sort, order string) {
	return "", "created_at", "desc"
}

// ==========================================
// FILTER HELPERS
// ==========================================

// FilterParams representa parâmetros de filtro
type FilterParams struct {
	Role   string `form:"role" binding:"omitempty,oneof=admin manager user guest" example:"user"`
	Status string `form:"status" binding:"omitempty,oneof=active inactive pending suspended" example:"active"`
}

// ==========================================
// RESPONSE BUILDERS
// ==========================================

// NewSuccessResponse cria uma resposta de sucesso
func NewSuccessResponse(message string, data interface{}) SuccessResponse {
	return SuccessResponse{
		Success:   true,
		Message:   message,
		Data:      data,
		Timestamp: time.Now(),
	}
}

// NewSuccessResponseWithRequestID cria uma resposta de sucesso com Request ID
func NewSuccessResponseWithRequestID(message string, data interface{}, requestID string) SuccessResponse {
	return SuccessResponse{
		Success:   true,
		Message:   message,
		Data:      data,
		RequestID: requestID,
		Timestamp: time.Now(),
	}
}

// NewErrorResponse cria uma resposta de erro
func NewErrorResponse(error, message string) ErrorResponse {
	return ErrorResponse{
		Success:   false,
		Error:     error,
		Message:   message,
		Timestamp: time.Now(),
	}
}

// NewErrorResponseWithRequestID cria uma resposta de erro com Request ID
func NewErrorResponseWithRequestID(error, message, requestID string) ErrorResponse {
	return ErrorResponse{
		Success:   false,
		Error:     error,
		Message:   message,
		RequestID: requestID,
		Timestamp: time.Now(),
	}
}

// NewErrorResponseWithTrace cria uma resposta de erro com Request ID e Trace ID
func NewErrorResponseWithTrace(error, message, requestID, traceID string) ErrorResponse {
	return ErrorResponse{
		Success:   false,
		Error:     error,
		Message:   message,
		RequestID: requestID,
		TraceID:   traceID,
		Timestamp: time.Now(),
	}
}

// NewValidationErrorResponse cria uma resposta de erro de validação
func NewValidationErrorResponse(details []ValidationErrorDetail) ValidationErrorResponse {
	return ValidationErrorResponse{
		Success: false,
		Error:   "Validation failed",
		Message: "One or more fields are invalid",
		Details: details,
	}
}

// ==========================================
// HTTP STATUS HELPERS
// ==========================================

// HTTPStatus representa códigos de status HTTP comuns
type HTTPStatus int

const (
	StatusOK                  HTTPStatus = 200
	StatusCreated             HTTPStatus = 201
	StatusNoContent           HTTPStatus = 204
	StatusBadRequest          HTTPStatus = 400
	StatusUnauthorized        HTTPStatus = 401
	StatusForbidden           HTTPStatus = 403
	StatusNotFound            HTTPStatus = 404
	StatusConflict            HTTPStatus = 409
	StatusUnprocessableEntity HTTPStatus = 422
	StatusInternalServerError HTTPStatus = 500
)

// String retorna a representação string do status
func (s HTTPStatus) String() string {
	switch s {
	case StatusOK:
		return "200 OK"
	case StatusCreated:
		return "201 Created"
	case StatusNoContent:
		return "204 No Content"
	case StatusBadRequest:
		return "400 Bad Request"
	case StatusUnauthorized:
		return "401 Unauthorized"
	case StatusForbidden:
		return "403 Forbidden"
	case StatusNotFound:
		return "404 Not Found"
	case StatusConflict:
		return "409 Conflict"
	case StatusUnprocessableEntity:
		return "422 Unprocessable Entity"
	case StatusInternalServerError:
		return "500 Internal Server Error"
	default:
		return "500 Internal Server Error"
	}
}

// ==========================================
// CONTEXT HELPERS
// ==========================================

// RequestContext representa informações do contexto da requisição
type RequestContext struct {
	UserID    string
	UserRole  string
	UserEmail string
	RequestID string
	IPAddress string
	UserAgent string
	Timestamp time.Time
}

// ==========================================
// METADATA HELPERS
// ==========================================

// ResponseMetadata representa metadados da resposta
type ResponseMetadata struct {
	RequestID   string    `json:"request_id,omitempty"`
	Timestamp   time.Time `json:"timestamp"`
	Version     string    `json:"version,omitempty"`
	Environment string    `json:"environment,omitempty"`
}

// AddMetadata adiciona metadados à resposta
func AddMetadata(data interface{}, requestID, version, environment string) map[string]interface{} {
	return map[string]interface{}{
		"data": data,
		"metadata": ResponseMetadata{
			RequestID:   requestID,
			Timestamp:   time.Now(),
			Version:     version,
			Environment: environment,
		},
	}
}
