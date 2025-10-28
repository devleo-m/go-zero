package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Data    interface{} `json:"data,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
	Message string      `json:"message,omitempty"`
	Error   string      `json:"error,omitempty"`
	Success bool        `json:"success"`
}

type Meta struct {
	Page       int   `json:"page,omitempty"`
	Limit      int   `json:"limit,omitempty"`
	Total      int64 `json:"total,omitempty"`
	TotalPages int   `json:"total_pages,omitempty"`
}

// Success retorna uma resposta de sucesso.
func Success(c *gin.Context, data interface{}, message ...string) {
	msg := ""
	if len(message) > 0 {
		msg = message[0]
	}

	c.JSON(http.StatusOK, Response{
		Success: true,
		Message: msg,
		Data:    data,
	})
}

// Created retorna uma resposta de criação bem-sucedida.
func Created(c *gin.Context, data interface{}, message ...string) {
	msg := "Created successfully"
	if len(message) > 0 {
		msg = message[0]
	}

	c.JSON(http.StatusCreated, Response{
		Success: true,
		Message: msg,
		Data:    data,
	})
}

// NoContent retorna uma resposta sem conteúdo.
func NoContent(c *gin.Context, message ...string) {
	msg := "Operation completed successfully"
	if len(message) > 0 {
		msg = message[0]
	}

	c.JSON(http.StatusNoContent, Response{
		Success: true,
		Message: msg,
	})
}

// Error retorna uma resposta de erro.
func Error(c *gin.Context, statusCode int, errorCode, message string) {
	c.JSON(statusCode, Response{
		Success: false,
		Error:   errorCode,
		Message: message,
	})
}

// BadRequest retorna uma resposta de erro de requisição inválida.
func BadRequest(c *gin.Context, errorCode, message string) {
	Error(c, http.StatusBadRequest, errorCode, message)
}

// Unauthorized retorna uma resposta de não autorizado.
func Unauthorized(c *gin.Context, errorCode, message string) {
	Error(c, http.StatusUnauthorized, errorCode, message)
}

// Forbidden retorna uma resposta de proibido.
func Forbidden(c *gin.Context, errorCode, message string) {
	Error(c, http.StatusForbidden, errorCode, message)
}

// NotFound retorna uma resposta de não encontrado.
func NotFound(c *gin.Context, errorCode, message string) {
	Error(c, http.StatusNotFound, errorCode, message)
}

// Conflict retorna uma resposta de conflito.
func Conflict(c *gin.Context, errorCode, message string) {
	Error(c, http.StatusConflict, errorCode, message)
}

// InternalServerError retorna uma resposta de erro interno do servidor.
func InternalServerError(c *gin.Context, errorCode, message string) {
	Error(c, http.StatusInternalServerError, errorCode, message)
}

// Paginated retorna uma resposta paginada.
func Paginated(c *gin.Context, data interface{}, meta *Meta, message ...string) {
	msg := ""
	if len(message) > 0 {
		msg = message[0]
	}

	c.JSON(http.StatusOK, Response{
		Success: true,
		Message: msg,
		Data:    data,
		Meta:    meta,
	})
}

// ValidationError retorna uma resposta de erro de validação.
func ValidationError(c *gin.Context, errors map[string]string) {
	c.JSON(http.StatusBadRequest, Response{
		Success: false,
		Error:   "VALIDATION_ERROR",
		Message: "Validation failed",
		Data:    gin.H{"errors": errors},
	})
}

// NewMeta cria uma nova estrutura de meta para paginação.
func NewMeta(page, limit int, total int64) *Meta {
	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	return &Meta{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	}
}
