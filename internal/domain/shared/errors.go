package shared

import (
	"fmt"
	"strings"
)

// DomainError representa um erro específico do domínio
// Diferente de erros de infraestrutura (DB, HTTP), estes são erros de negócio
type DomainError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// Error implementa a interface error do Go
func (e DomainError) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("[%s] %s: %s", e.Code, e.Message, e.Details)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// NewDomainError cria um novo erro de domínio
func NewDomainError(code, message string) DomainError {
	return DomainError{
		Code:    code,
		Message: message,
	}
}

// WithDetails adiciona detalhes ao erro
func (e DomainError) WithDetails(details string) DomainError {
	e.Details = details
	return e
}

// ValidationError representa erros de validação
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Value   string `json:"value,omitempty"`
}

// Error implementa a interface error
func (e ValidationError) Error() string {
	return fmt.Sprintf("validation error on field '%s': %s", e.Field, e.Message)
}

// ValidationErrors é uma coleção de erros de validação
type ValidationErrors []ValidationError

// Error implementa a interface error
func (ve ValidationErrors) Error() string {
	var messages []string
	for _, err := range ve {
		messages = append(messages, err.Error())
	}
	return strings.Join(messages, "; ")
}

// HasErrors verifica se há erros de validação
func (ve ValidationErrors) HasErrors() bool {
	return len(ve) > 0
}

// Add adiciona um erro de validação
func (ve *ValidationErrors) Add(field, message, value string) {
	*ve = append(*ve, ValidationError{
		Field:   field,
		Message: message,
		Value:   value,
	})
}
